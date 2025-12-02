// main.go
package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Book 结构体 —— 与 books 表字段一一对应
type Book struct {
	ID     int     `db:"id"`     // 主键
	Title  string  `db:"title"`  // 书名
	Author string  `db:"author"` // 作者
	Price  float64 `db:"price"`  // 价格（使用 float64 更安全）
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

	// ========== 作业要求：复杂查询 + 类型安全映射 ==========
	var expensiveBooks []Book

	// 复杂查询：价格 > 50 且按价格降序排列
	query := `
        SELECT id, title, author, price 
        FROM books 
        WHERE price > ? 
        ORDER BY price DESC`

	err = db.Select(&expensiveBooks, query, 50.0)
	if err != nil {
		log.Fatal("查询价格 > 50 的书籍失败:", err)
	}

	// 输出结果
	fmt.Printf("\n【价格大于 50 元的书籍（共 %d 本）】\n", len(expensiveBooks))
	for _, book := range expensiveBooks {
		fmt.Printf("ID: %d | 书名: 《%s》 | 作者: %s | 价格: ¥%.2f\n",
			book.ID, book.Title, book.Author, book.Price)
	}
}
