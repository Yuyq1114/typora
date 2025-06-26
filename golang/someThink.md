# 同时并发流水线使用：

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

