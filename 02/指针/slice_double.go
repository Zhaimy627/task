package main

import "fmt"

// doubleSlice 将切片中每个元素乘以 2
func doubleSlice(slicePtr *[]int) {
	// 解引用得到切片
	for i := range *slicePtr {
		(*slicePtr)[i] *= 2 // 修改原切片
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("修改前:", nums) // [1 2 3 4 5]

	doubleSlice(&nums) // 传入切片地址

	fmt.Println("修改后:", nums) // [2 4 6 8 10]
}
