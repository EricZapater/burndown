package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.Create(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func(h *Handler) GetUser(c *gin.Context) {
	user, err := h.service.FindByID(c.Request.Context(), uuid.MustParse(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func(h *Handler) ChangePassword(c *gin.Context) {
	var password string
	if err := c.ShouldBindJSON(&password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	err := h.service.ChangePassword(c.Request.Context(), uuid.MustParse(c.Param("id")), password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func(h *Handler) InactiveUser(c *gin.Context) {
	err := h.service.InactiveUser(c.Request.Context(), uuid.MustParse(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User inactivated successfully"})
}

func(h *Handler) AddMacros(c *gin.Context) {
	var macros Macros
	if err := c.ShouldBindJSON(&macros); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	err := h.service.AddMacros(c.Request.Context(), uuid.MustParse(c.Param("id")), macros)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Macros added successfully"})
}

func(h *Handler) GetMacros(c *gin.Context) {
	macros, err := h.service.GetMacros(c.Request.Context(), uuid.MustParse(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, macros)
}

func(h *Handler) GetHistoricalMacros(c *gin.Context) {
	macros, err := h.service.GetHistoricalMacros(c.Request.Context(), uuid.MustParse(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, macros)
}

func(h *Handler) DeleteMacros(c *gin.Context) {
	err := h.service.DeleteMacros(c.Request.Context(), uuid.MustParse(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Macros deleted successfully"})
}

func(h *Handler) UpdateMacros(c *gin.Context) {
	var macros Macros
	if err := c.ShouldBindJSON(&macros); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	err := h.service.UpdateMacros(c.Request.Context(), uuid.MustParse(c.Param("id")), macros)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Macros updated successfully"})
}

func(h *Handler) InactiveMacros(c *gin.Context) {
	err := h.service.InactiveMacros(c.Request.Context(), uuid.MustParse(c.Param("id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Macros inactivated successfully"})
}
