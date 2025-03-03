package main

import (
	connection "github.com/Dnreikronos/image_resizer_b/db"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := connection.OpenConnection()
	if err != nil {
		panic(err)
	}
}
