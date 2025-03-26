package main

import (
	"frappuccino/router"
	"log"
)

func main() {
	router.SetupRouter()
	log.Println("Server is running on http://0.0.0.0:8080")
}
