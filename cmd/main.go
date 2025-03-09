package main

import (
	"fmt"

	"github.com/Dnreikronos/image_resizer_b/configs"
	connection "github.com/Dnreikronos/image_resizer_b/db"
	h "github.com/Dnreikronos/image_resizer_b/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = configs.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	db, err := connection.OpenConnection()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/upload", h.UploadImage)
	r.GET("/image/:id", h.GetImage)
}
