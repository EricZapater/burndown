package auth

import (
	"api/internal/users"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Handler struct {
    service    Service
    jwtMiddleware *jwt.GinJWTMiddleware
}

func NewHandler(service Service, jwtMiddleware *jwt.GinJWTMiddleware) *Handler {
    return &Handler{
        service:    service,
        jwtMiddleware: jwtMiddleware,
    }
}

// Login processa una petició de login i retorna un token JWT
func (h *Handler) Login(c *gin.Context) {
    var loginRequest LoginRequest
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login request"})
        return
    }
    
    response, err := h.service.Login(c.Request.Context(), &loginRequest)
    if err != nil {
        var statusCode int
        switch err {
        case ErrInvalidCredentials:
            statusCode = http.StatusUnauthorized
        case users.ErrInactiveUser:
            statusCode = http.StatusForbidden
        default:
            statusCode = http.StatusInternalServerError
        }
        
        c.JSON(statusCode, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, LoginResponse{
        Token:  response.Token,
        Expire: response.Expire,
        User:   response.User,
    })
}

// RefreshToken utilitza el middleware JWT per refrescar el token
func (h *Handler) RefreshToken(c *gin.Context) {
    h.jwtMiddleware.RefreshHandler(c)
}
