package migration

import (
	"github.com/Dnreikronos/image_resizer_b/models"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	createTables(db)
}

func createTables(db *gorm.DB) {
	db.AutoMigrate(&models.Image{})
}
