package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codahale/sss"
	ecies "github.com/ecies/go/v2"
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

		case "dev-ecies":
			// 開発用
			k, err := ecies.GenerateKey()
			if err != nil {
				panic(err)
			}
			log.Println("key pair has been generated")

			ciphertext, err := ecies.Encrypt(k.PublicKey, []byte("THIS IS THE TEST"))
			if err != nil {
				panic(err)
			}
			log.Printf("plaintext encrypted: %v\n", ciphertext)

			plaintext, err := ecies.Decrypt(k, ciphertext)
			if err != nil {
				panic(err)
			}
			log.Printf("ciphertext decrypted: %s\n", string(plaintext))

		case "dev-se":
			// 開発用
			// 分散したシェアを更に楕円曲線暗号にて暗号化し、通信の際に暗号化されたシェアを送信する
			// 送信されたシェアを復号化し、各デバイスで保存する。(デバイスへの転送は別途実装する)
			// 複合時に再度結合する

			// 分散シェアの生成
			n := 5
			k := 3
			secret := "secret"

			privatekeys, shares := gr.Encrypt(n, k, secret, nil)

			// 分散シェアの復号化
			plain_shares := map[byte][]byte{}
			for i := 0; i < n; i++ {
				plain_shares[byte(i+1)], err = ecies.Decrypt(privatekeys[i], shares[i])
				if err != nil {
					panic(err)
				}
				log.Println("ciphertext decrypted: ", plain_shares[byte(i+1)])
			}

			// 分散シェアの結合
			recov := sss.Combine(plain_shares)
			fmt.Println(string(recov))

		default:
			log.Fatal("Error: No mode selected")
	}
}
