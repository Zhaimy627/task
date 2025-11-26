// main.go
package main

import (
	"github.com/gin-gonic/gin"
	"go-blog-system/config"
	"go-blog-system/handlers"
	"go-blog-system/middleware"
	"log"
)

func main() {
	config.InitDB()

	r := gin.Default()

	// 公开接口
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/posts", handlers.GetPosts)
	r.GET("/posts/:id", handlers.GetPost)
	r.GET("/posts/:id/comments", handlers.GetComments)

	// 认证接口
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/posts", handlers.CreatePost)
		api.PUT("/posts/:id", handlers.UpdatePost)
		api.DELETE("/posts/:id", handlers.DeletePost)
		api.POST("/posts/:id/comments", handlers.CreateComment)
	}

	log.Println("服务器启动: http://localhost:8081")
	r.Run(":8081")
}
