// handlers/comment.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"go-blog-system/config"
	"go-blog-system/models"
	"net/http"
	"strconv"
)

func CreateComment(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(models.User)
	postID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		Content: input.Content,
		UserID:  u.ID,
		PostID:  uint(postID),
	}

	config.DB.Create(&comment)
	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	var comments []models.Comment
	config.DB.Preload("User").Where("post_id = ?", postID).Find(&comments)
	c.JSON(http.StatusOK, comments)
}
