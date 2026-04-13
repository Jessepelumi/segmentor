package handlers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
		return
	}

	videoID := fmt.Sprintf("%d", time.Now().UnixNano())
	ext := filepath.Ext(file.Filename)
	savePath := filepath.Join("storage/uploads", videoID+ext)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create HLS output directory
	hlsDir := filepath.Join("storage/hls", videoID)
	if err := os.MkdirAll(hlsDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HLS directory"})
		return
	}

	// Run FFmpeg — blocking for now, background worker later
	if err := processVideoToHLS(savePath, hlsDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "FFmpeg processing failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Upload and processing successful",
		"video_id": videoID,
		"hls_path": hlsDir,
		"playlist": filepath.Join(hlsDir, "index.m3u8"),
	})
}

func processVideoToHLS(inputPath string, outputDir string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-codec:", "copy",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		outputDir+"/index.m3u8",
	)

	// Capture stderr so FFmpeg errors are visible in your Go logs
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
