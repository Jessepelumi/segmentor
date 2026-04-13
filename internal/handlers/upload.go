package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No video file provided",
		})
		return
	}

	// Create a simple unique ID (we'll improve later)
	videoID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Preserve extension
	ext := filepath.Ext(file.Filename)

	// Build path
	savePath := filepath.Join("storage/uploads", videoID+ext)

	// Save file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Upload successful",
		"video_id": videoID,
		"path":     savePath,
	})
}