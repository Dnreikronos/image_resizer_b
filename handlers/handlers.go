package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/Dnreikronos/image_resizer_b/db/connection"
	"github.com/Dnreikronos/image_resizer_b/models"
	"github.com/Dnreikronos/image_resizer_b/utils"
	"github.com/gin-gonic/gin"
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
