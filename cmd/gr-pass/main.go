package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main(){

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENV")
	fmt.Println("MODE:" + env)


	switch env {
		case "terminal":
			Terminal()
		default:
			log.Fatal("Error: No mode selected")
	}
}
