package main

import (
	"fmt"
	"time"
)

// printOdds 打印 1,3,5,7,9
func printOdds() {
	for i := 1; i <= 10; i += 2 {
		fmt.Printf("奇数协程: %d\n", i)
		time.Sleep(100 * time.Millisecond) // 模拟工作，观察交错输出
	}
}

// printEvens 打印 2,4,6,8,10
func printEvens() {
	for i := 2; i <= 10; i += 2 {
		fmt.Printf("偶数协程: %d\n", i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println("启动两个协程...")

	go printOdds()  // 启动奇数协程
	go printEvens() // 启动偶数协程

	// 主协程等待子协程完成（否则程序立即退出）
	time.Sleep(2 * time.Second)
	fmt.Println("主协程结束")
}
