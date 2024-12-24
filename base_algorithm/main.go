package main

import (
	"algorithm/basealgorithm"
	"fmt"
)

func main() {
	weights := []int{2, 3, 4, 5}
	values := []int{3, 4, 5, 6}
	capacity := 8

	fmt.Println(basealgorithm.KnapsackD(weights, values, capacity))
}
