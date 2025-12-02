package main

import "fmt"

// addTen 修改指针指向的值，增加 10
func addTen(ptr *int) {
	*ptr += 10 // 解引用后修改原值
}

func main() {
	num := 5
	fmt.Println("修改前:", num) // 输出: 5

	addTen(&num) // 传地址

	fmt.Println("修改后:", num) // 输出: 15
}
