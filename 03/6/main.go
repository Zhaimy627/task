// main.go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ==================== 模型定义（与上一题一致） ====================

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"type:varchar(50);unique;not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"type:varchar(200);not null"`
	Content   string `gorm:"type:text;not null"`
	UserID    uint
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Author   User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	PostID    uint
	UserID    uint
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Post   Post `gorm:"foreignKey:PostID"`
	Author User `gorm:"foreignKey:UserID"`
}

// ==================== 主函数 ====================
func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接失败: " + err.Error())
	}

	// 1. 自动建表
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 2. 清空旧数据 + 插入测试数据
	db.Exec("TRUNCATE TABLE comments")
	db.Exec("TRUNCATE TABLE posts")
	db.Exec("TRUNCATE TABLE users")

	insertTestData(db)

	// ========== 作业要求 1：查询用户的所有文章 + 评论 ==========
	userID := uint(1)
	var userWithPosts User
	db.Preload("Posts.Comments").First(&userWithPosts, userID)

	fmt.Printf("\n用户【%s】发布的所有文章及评论：\n", userWithPosts.Username)
	for _, post := range userWithPosts.Posts {
		fmt.Printf("  文章：《%s》\n", post.Title)
		if len(post.Comments) == 0 {
			fmt.Println("    （暂无评论）")
		} else {
			for _, c := range post.Comments {
				fmt.Printf("    评论：%s\n", c.Content)
			}
		}
	}

	// ========== 作业要求 2：查询评论最多的文章 ==========
	var topPost Post
	db.Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Group("posts.id").
		Order("COUNT(comments.id) DESC").
		First(&topPost)

	// 获取评论数量
	var commentCount int64
	db.Model(&Comment{}).Where("post_id = ?", topPost.ID).Count(&commentCount)

	fmt.Printf("\n评论最多的文章：\n")
	fmt.Printf("  标题：《%s》\n", topPost.Title)
	fmt.Printf("  作者：%s\n", getUsername(db, topPost.UserID))
	fmt.Printf("  评论数：%d 条\n", commentCount)
}

// ==================== 测试数据 ====================
func insertTestData(db *gorm.DB) {
	// 用户
	user1 := User{Username: "张三", Email: "zs@example.com", Password: "123"}
	user2 := User{Username: "李四", Email: "ls@example.com", Password: "123"}
	db.Create(&user1)
	db.Create(&user2)

	// 文章
	post1 := Post{Title: "GORM 真香", Content: "太强大了", UserID: user1.ID}
	post2 := Post{Title: "Go 语言入门", Content: "简单易学", UserID: user1.ID}
	post3 := Post{Title: "Python 还是香", Content: "对比一下", UserID: user2.ID}
	db.Create(&post1)
	db.Create(&post2)
	db.Create(&post3)

	// 评论
	db.Create(&Comment{Content: "同意！", PostID: post1.ID, UserID: user1.ID})
	db.Create(&Comment{Content: "我也是 GORM 粉", PostID: post1.ID, UserID: user2.ID})
	db.Create(&Comment{Content: "点赞", PostID: post1.ID, UserID: user1.ID})
	db.Create(&Comment{Content: "不错", PostID: post2.ID, UserID: user2.ID})
}

// ==================== 辅助函数 ====================
func getUsername(db *gorm.DB, userID uint) string {
	var user User
	db.Select("username").First(&user, userID)
	return user.Username
}
