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
