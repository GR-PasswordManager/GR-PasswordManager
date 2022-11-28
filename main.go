package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codahale/sss"
	"github.com/joho/godotenv"

	"github.com/GR-PasswordManager/GR-PasswordManager/cmd/gr"
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
			gr.Terminal()
		case "dev-sss":
			// 開発用
			n := 5
			k := 3
			secret := "secret"
			shares, err := sss.Split(byte(n), byte(k), []byte(secret))
			if err != nil {
				panic(err)
			}
			fmt.Println(shares)

			recov := sss.Combine(shares)
			fmt.Println(recov)
			fmt.Println(string(recov))

		default:
			log.Fatal("Error: No mode selected")
	}
}
