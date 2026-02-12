package advisor

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type AnalyzeRequest struct {
	Description string `form:"description"`
}

func (h *Handler) Analyze(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Get description
	description := c.PostForm("description")

	// Get image
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}
	defer file.Close()

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	// Read image data
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
		return
	}

	// Determine format (basic check based on extension or content type)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == ".jpg" || ext == ".jpeg" {
		ext = "jpeg"
	} else if ext == ".png" {
		ext = "png"
	} else {
		ext = "jpeg" 
	}
    
	// Call service
	retrials := 3
	for i := 0; i < retrials; i++ {
		result, err := h.service.AnalyzeImage(c.Request.Context(), fileBytes, ext, description)
		if err == nil {
			c.JSON(http.StatusOK, result)
			return
		}
		log.Println(i, " Attempt failed", err)
		time.Sleep(5 * time.Second)
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Analysis failed: %v", err)})
	return	
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/advisor/analyze", h.Analyze)
}
