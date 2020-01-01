package main

import (
	"log"
)
import "github.com/joho/godotenv"
import esi "github.com/dariusbakunas/eve-processors"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	esi.Process()
}