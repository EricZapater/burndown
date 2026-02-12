package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
// User management routes
	usersGroup := router.Group("/users")
	{
		usersGroup.GET("", handler.GetAll)
		usersGroup.GET("/:id", handler.GetUser)
		usersGroup.PUT("/:id/password", handler.ChangePassword)
		usersGroup.DELETE("/:id", handler.InactiveUser)
		
		// Macros routes
		usersGroup.POST("/:id/macros", handler.AddMacros)
		usersGroup.GET("/:id/macros", handler.GetMacros)
		usersGroup.GET("/:id/macros/historical", handler.GetHistoricalMacros)
		usersGroup.PUT("/:id/macros/:macroId", handler.UpdateMacros)
		usersGroup.DELETE("/:id/macros/:macroId", handler.DeleteMacros)
		usersGroup.PATCH("/:id/macros/:macroId/inactive", handler.InactiveMacros)
	}
}