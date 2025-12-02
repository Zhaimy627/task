package main

import (
	"fmt"
	"time"
)

func main() {
	// 创建无缓冲通道（容量为0）
	ch := make(chan int)

	// 生产者协程：生成 1~10 并发送
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("生产者发送: %d\n", i)
			ch <- i // 发送，阻塞直到被接收
		}
		close(ch) // 关闭通道
	}()

	// 消费者协程：接收并打印
	go func() {
		for num := range ch { // 自动遍历，直到通道关闭
			fmt.Printf("消费者接收: %d\n", num)
			time.Sleep(100 * time.Millisecond) // 模拟处理延迟
		}
		fmt.Println("通道已关闭，消费者退出")
	}()

	// 主协程等待
	time.Sleep(2 * time.Second)
	fmt.Println("主协程结束")
}
