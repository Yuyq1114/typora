package sort

// 归并排序算法（Merge Sort）：
// 将数组或列表分成两个子序列，
// 对每个子序列进行归并排序，然后将两个有序子序列合并成一个有序序列。

func merge(arr []int, left, mid, right int) {
	n1 := mid - left + 1
	n2 := right - mid

	// 创建临时数组，用来存储左半部分和右半部分的元素
	L := make([]int, n1)
	R := make([]int, n2)

	// 将左半部分的元素拷贝到临时数组 L 中
	for i := 0; i < n1; i++ {
		L[i] = arr[left+i]
	}

	// 将右半部分的元素拷贝到临时数组 R 中
	for j := 0; j < n2; j++ {
		R[j] = arr[mid+1+j]
	}

	// 将临时数组 L 和 R 合并到原数组 arr 中
	i := 0
	j := 0
	k := left

	for i < n1 && j < n2 {
		if L[i] <= R[j] {
			arr[k] = L[i]
			i++
		} else {
			arr[k] = R[j]
			j++
		}
		k++
	}

	for i < n1 {
		arr[k] = L[i]
		i++
		k++
	}

	for j < n2 {
		arr[k] = R[j]
		j++
		k++
	}
}

func MergeSort(arr []int, left, right int) {
	if left < right {
		// 找到中间点
		mid := (left + right) / 2

		// 分别对左半部分和右半部分进行排序
		MergeSort(arr, left, mid)
		MergeSort(arr, mid+1, right)

		// 合并两个有序序列
		merge(arr, left, mid, right)
	}
}

// func main() {
// 	arr := []int{4, 3, 2, 1, 5}
// 	fmt.Println("Unsorted array:", arr)

// 	MergeSort(arr, 0, len(arr)-1)
// 	fmt.Println("Sorted array:", arr)
// }
