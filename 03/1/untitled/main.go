package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:sa123456@tcp(127.0.0.1:3306)/ry?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}
	// 2. 自动迁移（如果表不存在就创建）
	db.AutoMigrate(&Student{})
	// 清空表（方便重复运行）
	db.Exec("TRUNCATE TABLE students")
	// ========== 作业：基本 CRUD ==========

	// 1. 插入：张三，20岁，三年级
	db.Create(&Student{Name: "张三", Age: 20, Grade: "三年级"})
	fmt.Println("插入张三成功")
	// 插入测试数据
	db.Create(&Student{Name: "李四", Age: 16, Grade: "二年级"})
	db.Create(&Student{Name: "王五", Age: 19, Grade: "三年级"})
	db.Create(&Student{Name: "赵六", Age: 14, Grade: "一年级"})
	// 2. 查询：年龄 > 18
	var adults []Student
	db.Where("age>?", 18).Find(&adults)
	fmt.Println("\n年龄 > 18 的学生：")
	for _, s := range adults {
		fmt.Printf("ID:%d | %s | %d岁 | %s\n", s.ID, s.Name, s.Age, s.Grade)
	}
	// 3. 更新：张三 → 四年级
	db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
	fmt.Println("张三已升入四年级")
	// 4. 删除：年龄 < 15
	db.Where("age < ?", 15).Delete(&Student{})
	fmt.Println("已删除未成年学生")

	// 最终结果
	var all []Student
	db.Find(&all)
	fmt.Println("\n最终 students 表数据：")
	for _, s := range all {
		fmt.Printf("ID:%d | %s | %d岁 | %s\n", s.ID, s.Name, s.Age, s.Grade)
	}
}
