package main

import (
	"fmt"

	"github.com/Dnreikronos/image_resizer_b/configs"
	connection "github.com/Dnreikronos/image_resizer_b/db"
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
}
