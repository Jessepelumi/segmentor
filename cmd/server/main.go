package main

import (
    "segmentor/internal/handlers"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.POST("/upload", handlers.UploadVideo)
    r.Static("/videos", "./storage/hls")
    r.StaticFile("/", "./client/index.html")

    r.Run(":8080")
}