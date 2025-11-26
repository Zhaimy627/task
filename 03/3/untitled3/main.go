// main.go
package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Employee 结构体（字段名与数据库列名一致，使用 db tag 映射）
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func main() {
	// 1. 连接数据库
	dsn := "root:sa123456@tcp(127.0.0.1:3306)/ry?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试连接
	if err = db.Ping(); err != nil {
		log.Fatal("Ping 失败:", err)
	}
	fmt.Println("数据库连接成功！")

	// ========== 作业要求 1：查询 "技术部" 员工 ==========
	var techEmployees []Employee
	query1 := `SELECT id, name, department, salary 
               FROM employees 
               WHERE department = ?`
	err = db.Select(&techEmployees, query1, "技术部")
	if err != nil {
		log.Fatal("查询技术部员工失败:", err)
	}

	fmt.Println("\n【技术部员工】")
	for _, emp := range techEmployees {
		fmt.Printf("ID: %d | 姓名: %s | 部门: %s | 工资: %.2f\n",
			emp.ID, emp.Name, emp.Department, emp.Salary)
	}

	// ========== 作业要求 2：查询工资最高员工 ==========
	var highestPaid Employee
	query2 := `SELECT id, name, department, salary 
               FROM employees 
               ORDER BY salary DESC 
               LIMIT 1`
	err = db.Get(&highestPaid, query2)
	if err != nil {
		log.Fatal("查询最高工资员工失败:", err)
	}

	fmt.Println("\n【工资最高员工】")
	fmt.Printf("ID: %d | 姓名: %s | 部门: %s | 工资: %.2f\n",
		highestPaid.ID, highestPaid.Name, highestPaid.Department, highestPaid.Salary)
}
