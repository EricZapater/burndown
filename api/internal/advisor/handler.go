package advisor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	jobID := uuid.New()
	requestMeta := json.RawMessage(`{}`) // Start with empty meta, can add headers/IP later if needed

	id, err := h.service.AnalyzeImage(c.Request.Context(), jobID, requestMeta, fileBytes, ext, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to submit analysis job: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": id,
		"status": "pending",
	})
}

func (h *Handler) GetAnalysisJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	job, err := h.service.GetAnalysisJob(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve job"})
		return
	}
	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/advisor/analyze", h.Analyze)
	r.GET("/advisor/jobs/:id", h.GetAnalysisJob)
}
