package sort

// 桶排序算法（Bucket Sort）：
// 将元素根据大小分配到不同的桶中，对每个桶内的元素进行排序，
// 然后按照桶的顺序依次输出所有元素，输出的顺序即为有序序列。
import (
	"math"
)

func BucketSort(arr []float64) []float64 {
	n := len(arr)
	if n <= 1 {
		return arr
	}
	max := arr[0]
	min := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	bucketSize := int(math.Floor((max-min)/(float64(n)-1))) + 1
	bucket := make([][]float64, bucketSize)
	for i := 0; i < len(bucket); i++ {
		bucket[i] = []float64{}
	}
	for _, v := range arr {
		index := int(math.Floor((v - min) / (float64(n) - 1) * float64(bucketSize-1)))
		bucket[index] = append(bucket[index], v)
	}
	result := []float64{}
	for i := 0; i < len(bucket); i++ {
		if len(bucket[i]) > 0 {
			bucket[i] = insertionSort(bucket[i])
			result = append(result, bucket[i]...)
		}
	}
	return result
}

// 桶内选择插入排序
func insertionSort(arr []float64) []float64 {
	for i := 1; i < len(arr); i++ {
		for j := i; j > 0; j-- {
			if arr[j] < arr[j-1] {
				arr[j], arr[j-1] = arr[j-1], arr[j]
			} else {
				break
			}
		}
	}
	return arr
}

// func main() {
//     arr := []float64{0.897, 0.565, 0.656, 0.1234, 0.665, 0.3434}
//     fmt.Println(BucketSort(arr))
// }
