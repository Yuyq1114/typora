package sort

// 计数排序算法（Counting Sort）：
// 统计每个元素出现的次数，然后依次输出每个元素，输出的顺序即为有序序列。

func CountingSort(arr []int) []int {
	n := len(arr)
	if n == 0 {
		return arr
	}
	min, max := arr[0], arr[0]
	for _, num := range arr {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}
	count := make([]int, max-min+1)
	for _, num := range arr {
		count[num-min]++
	}
	for i := 1; i < len(count); i++ {
		count[i] += count[i-1]
	}
	res := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		index := count[arr[i]-min] - 1
		res[index] = arr[i]
		count[arr[i]-min]--
	}
	return res
}

// func main() {
//     arr := []int{4, 10, 3, 5, 1}
//	   out := sort.CountingSort(arr)
//     fmt.Println(out)
// }
