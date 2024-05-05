package handler

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	r.GET("/", baseHandler)
	r.POST("/login", loginHandler)
	r.GET("/home", authMiddleware(), homeHandler)
	r.POST("/addBook", authMiddleware(), addBookHandler)
	r.DELETE("/deleteBook", authMiddleware(), deleteBookHandler)
}
