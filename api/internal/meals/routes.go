package meals

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	meals := router.Group("/meals")
	{
		meals.POST("", handler.Create)
		meals.GET("", handler.FindAll)
		meals.GET("/:id", handler.FindById)
		meals.PUT("/:id", handler.Update)
		meals.DELETE("/:id", handler.Delete)
	}

	// User-specific meal routes (Must use :id to match users module wildcard if registered on same group)
	users := router.Group("/users")
	{
		users.GET("/:id/meals", handler.FindAllByUserId)
		users.GET("/:id/meals/date/:date", handler.FindAllByUserIdAndDate)
	}
}
