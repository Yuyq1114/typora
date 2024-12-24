package sort

// 冒泡排序算法（Bubble Sort）：
// 依次比较相邻的两个元素，将较大的元素向后移动，
// 每次排序会将未排序部分的最大元素冒泡到已排序部分的末尾。
func BubbleSort(arr []int) []int {
	//外层循环遍历
	for i := 0; i < len(arr)-1; i++ {
		//内层循环0到前一个数
		for j := 0; j < len(arr)-i-1; j++ {
			//如果前一个数大于后面的数则交换
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

// func main() {
// 	arr := []int{4, 3, 2, 1, 5}
// 	fmt.Println("Unsorted array:", arr)

// 	arr = BubbleSort(arr)
// 	fmt.Println("Sorted array:", arr)
// }
