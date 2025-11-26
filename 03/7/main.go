// main.go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ==================== 模型定义 + 钩子函数 ====================

// User 用户（新增文章数量字段）
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"type:varchar(50);unique;not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	PostCount int       `gorm:"default:0"` // 文章数量统计
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Post 文章（新增评论状态字段）
type Post struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"type:varchar(200);not null"`
	Content       string `gorm:"type:text;not null"`
	UserID        uint
	CommentCount  int       `gorm:"default:0"` // 评论数量
	CommentStatus string    `gorm:"type:varchar(20);default:'无评论'"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	Author   User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

// Comment 评论
type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	PostID    uint
	UserID    uint
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Post   Post `gorm:"foreignKey:PostID"`
	Author User `gorm:"foreignKey:UserID"`
}

// ==================== Post 钩子：创建时更新用户文章数 ====================
func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
}

// ==================== Comment 钩子：删除时检查文章评论状态 ====================
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count)

	status := "有评论"
	if count == 0 {
		status = "无评论"
	}

	return tx.Model(&Post{}).Where("id = ?", c.PostID).
		Updates(map[string]interface{}{
			"comment_count":  count,
			"comment_status": status,
		}).Error
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

	// 2. 清空旧数据
	db.Exec("TRUNCATE TABLE comments")
	db.Exec("TRUNCATE TABLE posts")
	db.Exec("TRUNCATE TABLE users")

	// 3. 测试钩子
	testHooks(db)

	// 4. 最终结果
	printFinalData(db)
}

// ==================== 测试钩子函数 ====================
func testHooks(db *gorm.DB) {
	// 创建用户
	user := User{Username: "张三", Email: "zs@example.com", Password: "123"}
	db.Create(&user)
	fmt.Printf("创建用户：%s, 初始文章数：%d\n", user.Username, user.PostCount)

	// 创建文章 → 触发 Post.AfterCreate
	post := Post{Title: "GORM 钩子真香", Content: "自动更新！", UserID: user.ID}
	db.Create(&post)
	fmt.Println("创建文章 → 用户 PostCount 应 +1")

	// 创建 3 条评论
	for i := 1; i <= 3; i++ {
		db.Create(&Comment{Content: fmt.Sprintf("评论 %d", i), PostID: post.ID, UserID: user.ID})
	}
	fmt.Println("添加 3 条评论 → 文章状态应为 '有评论'")

	// 删除所有评论 → 触发 Comment.AfterDelete
	db.Where("post_id = ?", post.ID).Delete(&Comment{})
	fmt.Println("删除所有评论 → 文章状态应变为 '无评论'")
}

// ==================== 输出最终数据 ====================
func printFinalData(db *gorm.DB) {
	var user User
	db.Preload("Posts").First(&user, 1)

	fmt.Println("\n最终数据：")
	fmt.Printf("用户：%s, 文章数：%d\n", user.Username, user.PostCount)
	for _, p := range user.Posts {
		fmt.Printf("  文章：《%s》 | 评论数：%d | 状态：%s\n",
			p.Title, p.CommentCount, p.CommentStatus)
	}
}
