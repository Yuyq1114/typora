# 同时并发流水线使用：

## 1、管道做池：通信共享内存

```go
package main

import (
    "fmt"
    "time"
)

type Sample struct {
    Id    int
    Flag  bool
    Count [4]int
}

func main() {
    //potalPre := make(chan *Sample, 10)
    sem := make(chan struct{}, 12) // 控制最多 5 个样本在跑
    potal := make(chan *Sample, 12*4)
    //potalEnd := make(chan *Sample, 2)
    go func() {
       //id := 0
       for i := 0; i < 100; i++ {
          sem <- struct{}{} // 占名额
          sample := Sample{
             Id:    i,
             Flag:  false,
             Count: [4]int{0, 0, 0, 0},
          }

          potal <- &sample
          fmt.Println(" 生产了一个新的产品")

       }
    }()

    go func() {
       for {
          sample := <-potal
          if sample.Count[0] != 1 {
             fmt.Println("func0 is working")
             sample.Count[0] = 1
             time.Sleep(time.Second)
             fmt.Println("func0 work finished")
             if sample.Count[0] == 1 && sample.Count[1] == 1 && sample.Count[2] == 1 && sample.Count[3] == 1 {
                sample.Flag = true
                potal <- sample
             } else {
                potal <- sample
             }

          } else {
             potal <- sample
          }

       }

    }()
    go func() {
       for {
          sample := <-potal
          if sample.Count[3] != 1 {
             fmt.Println("func3 is working")
             sample.Count[3] = 1
             time.Sleep(time.Second)
             fmt.Println("func3 work finished")
             if sample.Count[0] == 1 && sample.Count[1] == 1 && sample.Count[2] == 1 && sample.Count[3] == 1 {
                sample.Flag = true
                potal <- sample
             } else {
                potal <- sample
             }
          } else {
             potal <- sample
          }
       }

    }()
    go func() {
       for {
          sample := <-potal
          if sample.Count[1] != 1 {
             fmt.Println("func1 is working")
             sample.Count[1] = 1
             time.Sleep(time.Second)
             fmt.Println("func1 work finished")
             if sample.Count[0] == 1 && sample.Count[1] == 1 && sample.Count[2] == 1 && sample.Count[3] == 1 {
                sample.Flag = true
                potal <- sample
             } else {
                potal <- sample
             }
          } else {
             potal <- sample
          }
       }

    }()
    go func() {
       for {
          sample := <-potal
          if sample.Count[2] != 1 {
             fmt.Println("func2 is working")
             sample.Count[2] = 1
             time.Sleep(time.Second)
             fmt.Println("func2 work finished")
             if sample.Count[0] == 1 && sample.Count[1] == 1 && sample.Count[2] == 1 && sample.Count[3] == 1 {
                sample.Flag = true
                potal <- sample
             } else {
                potal <- sample
             }
          } else {
             potal <- sample
          }
       }

    }()

    go func() {
       for {
          sample := <-potal
          if sample.Flag {
             time.Sleep(time.Second * 3)
             fmt.Printf("this is sample %d\n", sample.Id)
             <-sem // 释放名额
          } else {
             potal <- sample
          }
       }
    }()
    time.Sleep(1 * time.Hour)
}
```

## 2、共享内存做池：共享内存加锁通信

```
package main

import (
	"fmt"
	"sync"
	"time"
)

type Sample struct {
	Id    int
	Count [4]int
	Flag  bool
	mu    sync.Mutex // 每个 Sample 加锁保护
}

func worker(samples []*Sample, index int, wg *sync.WaitGroup) {
	for {
		allDone := true
		for _, s := range samples {
			s.mu.Lock()
			if s.Count[index] == 0 {
				s.Count[index] = 1
				fmt.Printf("func%d is working on sample %d\n", index, s.Id)
				s.mu.Unlock()
				time.Sleep(100 * time.Millisecond) // 模拟处理
				fmt.Printf("func%d finished on sample %d\n", index, s.Id)
			} else {
				s.mu.Unlock()
			}
		}
		// 检查是否所有样本都完成该阶段
		for _, s := range samples {
			s.mu.Lock()
			if s.Count[index] == 0 {
				allDone = false
			}
			s.mu.Unlock()
		}
		if allDone {
			wg.Done()
			return
		}
	}
}

func outputWatcher(samples []*Sample) {
	for {
		for _, s := range samples {
			s.mu.Lock()
			if !s.Flag && s.Count[0] == 1 && s.Count[1] == 1 && s.Count[2] == 1 && s.Count[3] == 1 {
				s.Flag = true
				fmt.Printf("✅ sample %d is done\n", s.Id)
			}
			s.mu.Unlock()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	samples := []*Sample{}
	for i := 0; i < 10; i++ {
		samples = append(samples, &Sample{Id: i})
	}

	var wg sync.WaitGroup
	wg.Add(4) // 四个阶段

	for i := 0; i < 4; i++ {
		go worker(samples, i, &wg)
	}

	go outputWatcher(samples)

	wg.Wait()
	fmt.Println("All stages done!")
}

```

## 3、ants协程池

```
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Sample struct {
	Id    int
	Count [4]int
	Flag  bool
	mu    sync.Mutex
}

func processStage(samples []*Sample, stage int, pool *ants.PoolWithFunc, wg *sync.WaitGroup) {
	var innerWg sync.WaitGroup
	for _, s := range samples {
		innerWg.Add(1)
		_ = pool.Invoke(&taskPayload{
			sample: s,
			stage:  stage,
			wg:     &innerWg,
		})
	}
	innerWg.Wait()
	wg.Done() // 当前 stage 完成
}

type taskPayload struct {
	sample *Sample
	stage  int
	wg     *sync.WaitGroup
}

// 实际处理函数
func taskFunc(v interface{}) {
	payload := v.(*taskPayload)
	s := payload.sample
	stage := payload.stage

	s.mu.Lock()
	if s.Count[stage] == 0 {
		s.Count[stage] = 1
		fmt.Printf("func%d is working on sample %d\n", stage, s.Id)
		s.mu.Unlock()

		time.Sleep(100 * time.Millisecond)

		fmt.Printf("func%d finished on sample %d\n", stage, s.Id)
	} else {
		s.mu.Unlock()
	}
	payload.wg.Done()
}

func outputWatcher(samples []*Sample, allDone *sync.WaitGroup) {
	defer allDone.Done()
	for {
		allFinished := true
		for _, s := range samples {
			s.mu.Lock()
			if !s.Flag && s.Count[0] == 1 && s.Count[1] == 1 && s.Count[2] == 1 && s.Count[3] == 1 {
				s.Flag = true
				fmt.Printf("✅ sample %d is done\n", s.Id)
			}
			if !s.Flag {
				allFinished = false
			}
			s.mu.Unlock()
		}
		if allFinished {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// 创建样本
	samples := []*Sample{}
	for i := 0; i < 10; i++ {
		samples = append(samples, &Sample{Id: i})
	}

	// 创建 ants 协程池（最多 8 个并发）
	pool, _ := ants.NewPoolWithFunc(8, taskFunc)
	defer pool.Release()

	var stageWg sync.WaitGroup
	stageWg.Add(4) // 四个阶段

	for i := 0; i < 4; i++ {
		go processStage(samples, i, pool, &stageWg)
	}

	var outputWg sync.WaitGroup
	outputWg.Add(1)
	go outputWatcher(samples, &outputWg)

	stageWg.Wait()
	outputWg.Wait()

	fmt.Println("All samples are fully processed!")
}

```

