package tests

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dnreikronos/image_resizer_b/handlers"
	"github.com/Dnreikronos/image_resizer_b/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	err = db.Migrator().DropTable(&models.Image{})
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}
	if err := db.AutoMigrate(&models.Image{}); err != nil {
		panic("Failed to migrate database")
	}
	return db
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/upload", handlers.UploadImage)
	r.GET("/image/:id", handlers.GetImage)
	r.PUT("/resize", handlers.ResizeImage)
	r.GET("/download/:id", handlers.DownloadResizedImage)
	return r
}

func TestUploadImage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB()
	r := SetupRouter(db)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake_image_data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Image uploaded successfully")
}
