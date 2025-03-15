package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Dnreikronos/image_resizer_b/models"
	"github.com/Dnreikronos/image_resizer_b/utils"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getDBFromContext(c *gin.Context) (*gorm.DB, bool) {
	db, exists := c.Get("db")
	if !exists {
		log.Println("Database connection is missing in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is missing"})
		return nil, false
	}

	gormDB, ok := db.(*gorm.DB)
	if !ok {
		log.Println("Invalid database instance in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database instance"})
		return nil, false
	}

	return gormDB, true
}

func UploadImage(c *gin.Context) {
	db, ok := getDBFromContext(c)
	if !ok {
		return
	}

	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		log.Println("Error: File is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	if err := utils.ValidateFile(fileHeader); err != nil {
		log.Println("Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Failed to read file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	image := models.Image{Filename: fileHeader.Filename, Data: fileData}

	if err := db.Create(&image).Error; err != nil {
		log.Println("Failed to save image:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	log.Println("Image uploaded successfully:", fileHeader.Filename)
	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "id": image.ID})
}

func GetImage(c *gin.Context) {
	db, ok := getDBFromContext(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var image models.Image
	if err := db.Where("id = ?", id).First(&image).Error; err != nil {
		log.Println("Image not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.Header("Content-Type", "image/jpeg")
	c.Writer.Write(image.Data)
}

func ResizeImage(c *gin.Context) {
	db, ok := getDBFromContext(c)
	if !ok {
		return
	}

	id := c.Query("id")
	widthParam := c.Query("width")
	heightParam := c.Query("height")

	log.Printf("Received resize request - ID: %s, Width: %s, Height: %s", id, widthParam, heightParam)

	width, err := strconv.Atoi(widthParam)
	if err != nil || width <= 0 {
		log.Printf("Invalid width parameter: %s", widthParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid width"})
		return
	}

	height, err := strconv.Atoi(heightParam)
	if err != nil || height <= 0 {
		log.Printf("Invalid height parameter: %s", heightParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid height"})
		return
	}

	var originalImage models.Image
	if err := db.First(&originalImage, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Image not foud for ID: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		log.Printf("Database error while fetching image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("Original image found: %s, Size: %d bytes", originalImage.Filename, len(originalImage.Data))

	img, format, err := image.Decode(bytes.NewReader(originalImage.Data))
	if err != nil {
		log.Printf("Failed to decode image for ID: %s, Error: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	var buf bytes.Buffer
	if format == "png" {
		err = imaging.Encode(&buf, resizedImg, imaging.PNG)
	} else {
		err = imaging.Encode(&buf, resizedImg, imaging.JPEG)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode resized image"})
		return
	}

	resizedImage := models.Image{
		ID:       uuid.New(),
		Filename: fmt.Sprintf("resized_%dx%d_%s", width, height, originalImage.Filename),
		Data:     buf.Bytes(),
	}

	if err := db.Create(&resizedImage).Error; err != nil {
		log.Printf("Failed to save resized image to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resized image"})
		return
	}

	log.Printf("Image resized successfully: %s", resizedImage.Filename)

	c.JSON(http.StatusOK, gin.H{"message": "Image resized", "id": resizedImage.ID})
}

func DownloadResizedImage(c *gin.Context) {
	db, ok := getDBFromContext(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var image models.Image
	if err := db.First(&image, id).Error; err != nil {
		log.Println("Image not found:", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+image.Filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(image.Data)))

	c.Data(http.StatusOK, "application/octet-stream", image.Data)
}
