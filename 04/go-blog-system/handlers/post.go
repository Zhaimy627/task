// handlers/post.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"go-blog-system/config"
	"go-blog-system/models"
	"net/http"
	"strconv"
)

func CreatePost(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(models.User)

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  u.ID,
	}

	config.DB.Create(&post)
	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Preload("User").Find(&posts)
	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var post models.Post
	if err := config.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(models.User)
	id, _ := strconv.Atoi(c.Param("id"))

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if post.UserID != u.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&post).Updates(input)
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(models.User)
	id, _ := strconv.Atoi(c.Param("id"))

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if post.UserID != u.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	config.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
