package main

import connection "github.com/Dnreikronos/image_resizer_b/db"




func main () {

	db, err := connection.OpenConnection()
	if err != nil {
		panic(err)
	}
}
