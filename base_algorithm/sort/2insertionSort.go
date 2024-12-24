package sort

// 插入排序算法（Insertion Sort）：
// 将未排序的元素依次插入已排序的部分中，
// 每次插入会将已排序部分中大于该元素的元素后移。

func InsertionSort(arr []int) []int {
	//遍历数组，从1开始
	for i := 1; i < len(arr); i++ {
		//保存第i的数
		key := arr[i]
		//j是前一个数
		j := i - 1
		//其他元素后移，直到第一个元素或者找到小于key的元素停止
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		//插入待排序元素
		arr[j+1] = key
	}
	return arr
}

// func main() {
// 	arr := []int{4, 3, 2, 1, 5}
// 	fmt.Println("Unsorted array:", arr)

// 	arr = InsertionSort(arr)
// 	fmt.Println("Sorted array:", arr)
// }
