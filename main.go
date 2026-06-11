package main

import (
	"log"
)

func main() {
	if err := Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer DB.Close()
}
