package main

import (
	"go_ai/db"
	"log"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}
	startServer()
}
