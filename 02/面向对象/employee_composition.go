package main

import "fmt"

// Person 基础结构体
type Person struct {
	Name string
	Age  int
}

// Employee 组合 Person，增加 EmployeeID
type Employee struct {
	Person     // 匿名嵌入（组合）
	EmployeeID string
}

// PrintInfo 方法：接收者是 Employee
func (e Employee) PrintInfo() {
	fmt.Printf("员工ID: %s\n", e.EmployeeID)
	fmt.Printf("姓名: %s\n", e.Name)
	fmt.Printf("年龄: %d\n", e.Age)
}

func main() {
	// 创建 Employee 实例
	emp := Employee{
		Person:     Person{Name: "张三", Age: 30},
		EmployeeID: "E12345",
	}

	// 调用方法
	fmt.Println("=== 员工信息 ===")
	emp.PrintInfo()

	// 直接访问嵌入字段（Go 的“继承”效果）
	fmt.Printf("直接访问: %s 是 %d 岁\n", emp.Name, emp.Age)
}
