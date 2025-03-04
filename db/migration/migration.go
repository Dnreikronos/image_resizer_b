import (
	"github.com/Dnreikronos/image_resizer_b/models"
	"gorm.io/gorm"
)
func RunMigration(db *gorm.DB){
	createTables(db)
}
