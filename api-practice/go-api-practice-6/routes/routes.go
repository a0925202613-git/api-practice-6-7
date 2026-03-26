package routes

import (
	"net/http"

	"go-api-practice-6/handlers"
	"go-api-practice-6/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.Logger())

	// 公開：菜單查詢
	api.GET("/menus", handlers.GetMenus)
	api.GET("/menus/:id", handlers.GetMenuByID)
	// 公開：訂單列表（可加 ?status=pending|completed|cancelled）
	api.GET("/orders", handlers.GetOrders)
	api.GET("/orders/:id", handlers.GetOrderByID)

	// 需 token：菜單管理、下單、取消訂單
	protected := api.Group("").Use(middleware.TokenAuth())
	{
		protected.POST("/menus", handlers.CreateMenu)
		protected.PUT("/menus/:id", handlers.UpdateMenu)
		protected.DELETE("/menus/:id", handlers.DeleteMenu)
		protected.POST("/orders", handlers.CreateOrder)
		protected.PATCH("/orders/:id/cancel", handlers.CancelOrder)
		protected.PUT("/orders/:id", handlers.UpdateOrder)
	}

	api.GET("/me", middleware.TokenAuth(), func(c *gin.Context) {
		token, _ := c.Get("token")
		c.JSON(http.StatusOK, gin.H{
			"message": "你有帶正確的 token！",
			"token":   token,
		})
	})
}
