// main.go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ==================== 模型定义 ====================

// User 用户
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"type:varchar(50);unique;not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// 一对多：一个用户有多篇文章
	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Post 文章
type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"type:varchar(200);not null"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint      // 外键（自动关联 User.ID）
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// 所属用户（反向关系，可选）
	Author User `gorm:"foreignKey:UserID"`

	// 一对多：一篇文章有多个评论
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

// Comment 评论
type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	PostID    uint      // 外键（自动关联 Post.ID）
	UserID    uint      // 评论者
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// 所属文章和用户（反向关系）
	Post   Post `gorm:"foreignKey:PostID"`
	Author User `gorm:"foreignKey:UserID"`
}

// ==================== 主函数：自动建表 ====================
func main() {
	// 1. 连接数据库
	dsn := "root:sa123456@tcp(127.0.0.1:3306)/ry?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 2. 自动迁移：创建表 + 外键关系
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("建表失败: " + err.Error())
	}

	fmt.Println("所有表创建成功！")
	fmt.Println("表结构：users, posts, comments")
	fmt.Println("关系：User 1:N Post, Post 1:N Comment")

	// 可选：插入测试数据验证关系
	createTestData(db)
}

// ==================== 测试数据（可选） ====================
func createTestData(db *gorm.DB) {
	// 创建用户
	user := User{
		Username: "张三",
		Email:    "zhangsan@example.com",
		Password: "hashed_password",
	}
	db.Create(&user)

	// 创建文章
	post := Post{
		Title:   "GORM 真香",
		Content: "今天学习了 GORM 的外键关系，太强大了！",
		UserID:  user.ID,
	}
	db.Create(&post)

	// 创建评论
	comment := Comment{
		Content: "楼主讲得太好了！",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	db.Create(&comment)

	fmt.Println("\n测试数据插入成功！")
}
