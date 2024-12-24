package sort

// 快速排序算法（Quick Sort）：
// 选择一个基准元素，将数组或列表分成两部分，其中一部分所有元素小于基准元素，
// 另一部分所有元素大于等于基准元素，然后递归对两部分进行快速排序。

// 快速排序函数
func QuickSort(arr []int, left, right int) {
	if left < right {
		// 分割点索引
		index := partition(arr, left, right)
		// 递归排序分割点左右两部分
		QuickSort(arr, left, index-1)
		QuickSort(arr, index+1, right)
	}
}

// 分割函数
func partition(arr []int, left, right int) int {
	// 选定基准数，这里选择数组最后一个元素
	pivot := arr[right]
	// 定义左右指针
	i, j := left, right-1
	for i <= j {
		// 从左到右查找大于基准数的元素索引
		for i <= j && arr[i] <= pivot {
			i++
		}
		// 从右到左查找小于基准数的元素索引
		for i <= j && arr[j] >= pivot {
			j--
		}
		// 如果左右指针没有相遇，则交换左右指针对应的元素
		if i <= j {
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	// 将基准数与左右指针相遇处的元素交换
	arr[i], arr[right] = arr[right], arr[i]
	// 返回分割点索引
	return i
}

// func main() {
//     arr := []int{3, 9, 1, 4, 7, 2, 8, 5, 6}
//     quickSort(arr, 0, len(arr)-1)
//     fmt.Println(arr) // [1 2 3 4 5 6 7 8 9]
// }
