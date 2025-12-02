package main

import (
	"fmt"
	"math"
)

// Shape 接口：所有形状必须实现
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 矩形
type Rectangle struct {
	Width  float64
	Height float64
}

// Area 实现
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter 实现
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle 圆形
type Circle struct {
	Radius float64
}

// Area 实现
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter 实现
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main() {
	// 创建实例
	rect := Rectangle{Width: 3, Height: 4}
	circle := Circle{Radius: 5}

	// 多态调用
	shapes := []Shape{rect, circle}

	fmt.Println("=== 形状计算结果 ===")
	for _, s := range shapes {
		fmt.Printf("面积: %.2f, 周长: %.2f\n", s.Area(), s.Perimeter())
	}
}
