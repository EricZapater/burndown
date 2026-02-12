package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtMiddleware *jwt.GinJWTMiddleware) {
	router.POST("/login", handler.Login)
	router.GET("/refresh_token", jwtMiddleware.RefreshHandler)
}