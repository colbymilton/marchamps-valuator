package main

import (
	"log"

	"github.com/colbymilton/marchamps-valuator/internal/restserver"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Welcome to the Marvel Champions Pack Valuator!")

	godotenv.Load()

	restserver.Run()
}
