package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"geekedu/common/oss"

	"github.com/gin-gonic/gin"
)

// UploadFile upload file to OSS or local storage
func UploadFile(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("failed to get file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "failed to get file: " + err.Error(),
		})
		return
	}

	log.Printf("upload file: %s, size: %d bytes", file.Filename, file.Size)

	// Open file
	src, err := file.Open()
	if err != nil {
		log.Printf("failed to open file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to open file",
		})
		return
	}
	defer src.Close()

	// Read file content to memory
	fileData, err := io.ReadAll(src)
	if err != nil {
		log.Printf("failed to read file content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to read file content: " + err.Error(),
		})
		return
	}

	// OSS - Use global client if available
	ossClient := oss.GetDefaultClient()
	if ossClient != nil {
		log.Printf("using OSS upload with global client")

		fileURL, err := ossClient.UploadFile(bytes.NewReader(fileData), filepath.Base(file.Filename))
		if err != nil {
			log.Printf("failed to upload file to OSS: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "failed to upload file: " + err.Error(),
			})
			return
		}

		log.Printf("file uploaded successfully: %s", fileURL)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "upload success",
			"data": gin.H{
				"url": fileURL,
			},
		})
		return
	}

	// Use local file storage
	log.Printf("OSS config incomplete, using local storage")
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to create upload directory",
		})
		return
	}

	// Generate unique filename
	timestamp := time.Now().Format("20060102150405")
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s_%d%s", timestamp, time.Now().UnixNano()%10000, ext)
	dstPath := filepath.Join(uploadDir, newFilename)

	// Write file content to destination
	if err := os.WriteFile(dstPath, fileData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to save file: " + err.Error(),
		})
		return
	}

	// Return local file URL
	fileURL := fmt.Sprintf("/uploads/%s", newFilename)
	log.Printf("file saved successfully: %s", fileURL)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "upload success",
		"data": gin.H{
			"url": fileURL,
		},
	})
}

// SignOSSURL sign OSS URL
func SignOSSURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid params",
		})
		return
	}

	ossClient := oss.GetDefaultClient()
	if ossClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "OSS client not initialized",
		})
		return
	}

	signedURL, err := ossClient.GetPresignedURL(req.URL, 3600)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to generate signed URL: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"signedUrl": signedURL,
		},
	})
}
