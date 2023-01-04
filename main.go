package main

import (
	"crypto/aes"
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

		case "dev-shares":
			// 開発用
			// ターミナルからの入力を受け取り、AESによる暗号化し分散、各デバイスに転送可能な状態にする。
			// また、複合・結合を行い、ターミナルに出力する。

			// 分散シェアの生成
			k := 3
			n := 5
			secret := ""

			fmt.Print("Input secret:")
			fmt.Scan(&secret)

			keys := map[byte]*ecies.PrivateKey{}
			publickeys := map[byte]*ecies.PublicKey{}
			for i := 1; i <= n; i++ {
				keys[byte(i)], err = ecies.GenerateKey()
				if err != nil {
					panic(err)
				}
				publickeys[byte(i)] = keys[byte(i)].PublicKey
			}

			block, err := aes.NewCipher([]byte("12345678901234561234567890123456"))
			if err != nil {
				panic(err)
			}

			shares := gr.Encrypt(k, n, []byte(secret), block, publickeys)

			// 分散シェアの復号化
			plain_shares := map[byte][]byte{}
			for i := 1; i <= len(shares); i++ {
				plain_shares[byte(i)], err = ecies.Decrypt(keys[byte(i)], shares[byte(i)])
				if err != nil {
					panic(err)
				}
				log.Println("ciphertext decrypted: ", plain_shares[byte(i)])
			}

			decrypt_text := gr.Decrypt(plain_shares, block)
			log.Println(string(decrypt_text))

		default:
			log.Fatal("Error: No mode selected")
	}
}
