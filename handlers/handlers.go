package handlers

import (
	"bytes"
	"image"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Dnreikronos/image_resizer_b/db/connection"
	"github.com/Dnreikronos/image_resizer_b/models"
	"github.com/Dnreikronos/image_resizer_b/utils"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadImage(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	if err := utils.ValidateFile(fileHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	db, _ := connection.OpenConnection()
	image := models.Image{Filename: fileHeader.Filename, Data: fileData}
	if err := db.Create(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	log.Println("Image uploaded:", fileHeader.Filename)
	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded sucessfully", "id": image.ID})
}

func GetImage(c *gin.Context) {
	id := c.Param("id")

	db, _ := connection.OpenConnection()
	var image models.Image

	if err := db.First(&image, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.Header("Context-Type", "image/jpeg")
	c.Writer.Write(image.Data)
}

func ReziseImage(c *gin.Context) {
	db, _ := connection.OpenConnection()

	id := c.Query("id")
	widthParam := c.Query("width")
	heightParam := c.Query("height")

	// convert size of them to INT
	width, err := strconv.Atoi(widthParam)
	if err != nil || width <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Width"})
		return
	}
	height, err := strconv.Atoi(heightParam)
	if err != nil || height <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Height"})
		return
	}

	var originalImage models.Image
	if err := db.First(&originalImage, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	img, format, err := image.Decode(bytes.NewReader(originalImage.Data))
	if err != nil {
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
		Filename: "resized_" + originalImage.Filename,
		Data:     buf.Bytes(),
	}
	if err := db.Create(&resizedImage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resized image"})
		return
	}

	log.Println("Image resized:", resizedImage.Filename)
	c.JSON(http.StatusOK, gin.H{"message": "Image resized", "id": resizedImage.ID})
}


func DownloadResizedImage(c *gin.Context) {
	id := c.Param("id")

	db, _ := connection.OpenConnection()
	var image models.Image

	if err := db.First(&image, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+image.Filename)
	c.Header("Context-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(image.Data)))

	c.Data(http.StatusOK, "application/octet-stream", image.Data)
}
