package sort

// 选择排序算法（Selection Sort）：
// 每次从未排序部分中选择最小的元素，将其移动到已排序部分的末尾。

func SelectionSort(arr []int) []int {
	//遍历
	for i := 0; i < len(arr)-1; i++ {
		//设置位置信息
		minIdx := i
		//在剩余序列中找到最小的排到已排序末尾
		for j := i + 1; j < len(arr); j++ {
			//j位置值小于最先的值则交换位置
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		//交换
		arr[i], arr[minIdx] = arr[minIdx], arr[i]
	}
	return arr
}

// func main() {
// 	arr := []int{4, 3, 2, 1, 5}
// 	fmt.Println("Unsorted array:", arr)

// 	arr = SelectionSort(arr)
// 	fmt.Println("Sorted array:", arr)
// }
