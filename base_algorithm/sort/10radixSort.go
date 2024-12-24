package sort

// 基数排序算法（Radix Sort）：
// 按照元素的个位、十位、百位等依次进行排序，
// 每次排序会保持之前的排序结果，直到所有位都排序完成，最终输出有序序列。
func RadixSort(arr []int) []int {
	n := len(arr)
	if n <= 1 {
		return arr
	}
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	radix := 10
	exp := 1
	result := make([]int, n)
	for max/exp > 0 {
		count := make([]int, radix)
		for _, v := range arr {
			index := (v / exp) % radix
			count[index]++
		}
		for i := 1; i < radix; i++ {
			count[i] += count[i-1]
		}
		for i := n - 1; i >= 0; i-- {
			index := (arr[i] / exp) % radix
			result[count[index]-1] = arr[i]
			count[index]--
		}
		for i := 0; i < n; i++ {
			arr[i] = result[i]
		}
		exp *= 10
	}
	return arr
}

// func main() {
//     arr := []int{170, 45, 75, 90, 802, 24, 2, 66}
//     fmt.Println(RadixSort(arr))
// }
