// models/comment.go
package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `json:"user_id"`
	PostID  uint   `json:"post_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post    Post   `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
