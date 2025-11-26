// config/db.go
package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:sa123456@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 如果你用 SQL 建表，这里可以注释掉 AutoMigrate
	// db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	DB = db
	log.Println("数据库连接成功")
}
