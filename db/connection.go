package connection

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "resizer"
	dbPassword = "admin123"
	dbName     = "resizer"
	dbTimeZone = "UTC"
)

func OpenConnection() (*gorm.DB, error) {
	connectionInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		dbTimeZone)

	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Sucesfuly connected to the database!")
	return db, nil
}
