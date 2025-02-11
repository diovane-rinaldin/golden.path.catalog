package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	s3Client *s3.Client
}

func NewUploadHandler(s3Client *s3.Client) *UploadHandler {
	return &UploadHandler{s3Client: s3Client}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Upload to S3
	// ... S3 upload logic ...

	c.JSON(200, gin.H{
		"url": fmt.Sprintf("%s/%s", os.Getenv("S3_BUCKET_URL"), filename),
	})
}
