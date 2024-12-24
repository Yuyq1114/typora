package sort

// 希尔排序算法（Shell Sort）：
// 将数组或列表分成若干个子序列，对每个子序列进行插入排序，
// 然后逐步缩小子序列的长度，直到整个数组或列表有序。
func ShellSort(arr []int) []int {
	n := len(arr)
	//增量序列折半
	gap := n / 2
	for gap > 0 {
		for i := gap; i < n; i++ {
			temp := arr[i]
			j := i
			for j >= gap && arr[j-gap] > temp {
				arr[j] = arr[j-gap]
				j -= gap
			}
			arr[j] = temp
		}
		gap /= 2
	}
	return arr
}

// func main() {
// 	arr := []int{4, 3, 2, 1, 5}
// 	fmt.Println("Unsorted array:", arr)

// 	arr = ShellSort(arr)
// 	fmt.Println("Sorted array:", arr)
// }
