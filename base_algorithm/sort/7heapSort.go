package sort

// 堆排序算法（Heap Sort）：
// 将数组或列表构建成一个二叉堆，
// 每次取出堆顶元素，并调整剩余元素使其重新构成一个二叉堆，
// 直到所有元素都被取出。

func HeapSort(arr []int) []int {
	n := len(arr)
	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}
	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		heapify(arr, i, 0)
	}
	return arr
}

func heapify(arr []int, n int, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2
	if left < n && arr[left] > arr[largest] {
		largest = left
	}
	if right < n && arr[right] > arr[largest] {
		largest = right
	}
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

// func main() {
//     arr := []int{4, 10, 3, 5, 1}
//     fmt.Println(heapSort(arr))
// }
