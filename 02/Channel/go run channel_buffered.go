package main

import (
	"fmt"
	"time"
)

func main() {
	// 创建缓冲通道，容量为 10
	ch := make(chan int, 10)

	// 生产者协程：发送 1~100
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i // 发送到缓冲区，若满则阻塞
			fmt.Printf("生产者发送: %d (缓冲区剩余: %d)\n", i, len(ch))
		}
		close(ch)
		fmt.Println("生产者完成，关闭通道")
	}()

	// 消费者协程：接收并打印
	go func() {
		for num := range ch {
			fmt.Printf("消费者接收: %d\n", num)
			time.Sleep(50 * time.Millisecond) // 模拟慢消费
		}
		fmt.Println("消费者完成")
	}()

	// 主协程等待
	time.Sleep(8 * time.Second)
	fmt.Println("主协程结束")
}
