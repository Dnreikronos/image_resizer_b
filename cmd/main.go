package main

import (
	"fmt"

	"github.com/Dnreikronos/image_resizer_b/configs"

	"github.com/Dnreikronos/image_resizer_b/db/connection"
	"github.com/Dnreikronos/image_resizer_b/db/migration"
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

	migration.RunMigration(db)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/upload", h.UploadImage)
	r.GET("/image/:id", h.GetImage)
	r.POST("/resize", h.ReziseImage)
}
