package routes

import (
	"net/http"

	"go-api-practice-7/handlers"
	"go-api-practice-7/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.Logger())

	// 公開：書籍列表（可加 ?available=true|false）、單筆、借閱列表、單筆借閱
	api.GET("/books", handlers.GetBooks)
	api.GET("/books/:id", handlers.GetBookByID)
	api.GET("/borrowals", handlers.GetBorrowals)
	api.GET("/borrowals/:id", handlers.GetBorrowalByID)

	// 需 token：書籍管理、借書、還書
	protected := api.Group("").Use(middleware.TokenAuth())
	{
		protected.POST("/books", handlers.CreateBook)
		protected.PUT("/books/:id", handlers.UpdateBook)
		protected.DELETE("/books/:id", handlers.DeleteBook)
		protected.POST("/borrowals", handlers.CreateBorrowal)
		protected.POST("/borrowals/:id/return", handlers.ReturnBorrowal)
	}

	api.GET("/me", middleware.TokenAuth(), func(c *gin.Context) {
		token, _ := c.Get("token")
		c.JSON(http.StatusOK, gin.H{
			"message": "你有帶正確的 token！",
			"token":   token,
		})
	})
}
