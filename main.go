package main

import (
	"bufio"
	"crypto/aes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

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

		case "dev-presentation":
			// プレゼンテーション用
			// ターミナルからの入力を受け取り、AESによる暗号化し分散、各デバイスに転送可能な状態にする。
			// ディレクトリを作成し、デバイスに転送する予定のファイルを作成する。
			// 一時停止
			// また、複合・結合を行い、ターミナルに出力する。

			// 分散シェアの生成
			k := 0;
			n := 0;

			fmt.Print("Input k:")
			fmt.Scan(&k)
			fmt.Print("Input n:")
			fmt.Scan(&n)

			if k > n {
				log.Fatal("Error: k > n")
			}else if k < 1 {
				log.Fatal("Error: k < 1")
			}else if n < 1 {
				log.Fatal("Error: n < 1")
			}

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

			log.Println(shares)

			// ディレクトリ作成
			dir := "presentation_" + time.Now().Format("2006-1-2_15-04")
			err = os.Mkdir(dir, 0777)

			// 未接続のドングル指定
			dongle := map[byte]string{}
			fmt.Println("未接続のドングルを指定してください。")
			for i := 1; i <= n; i++ {
				fmt.Println("ドングル" + strconv.Itoa(i) + "は接続していますか？(y/n)")
				var yn string
				fmt.Scan(&yn)
				if yn == "y" || yn == "n" {
					dongle[byte(i)] = string(yn)
				} else {
					log.Fatal("Error: y or n")
				}
			}

			// ディレクトリ作成
			for i := 1; i <= n; i++ {
				if string(dongle[byte(i)]) == "y" {
					err = os.Mkdir(dir + "/" + strconv.Itoa(i), 0777)
				}
			}

			// ファイル作成
			for i := 1; i <= len(shares); i++ {
				if f, err := os.Stat(dir + "/" + strconv.Itoa(i)); os.IsNotExist(err) || !f.IsDir() {
					fmt.Println("ドングル" + strconv.Itoa(i) + "は接続していません。")
					for j := 1; j <= n; j++ {
						if string(dongle[byte(j)]) == "y" {
							file, err := os.Create(dir + "/" + strconv.Itoa(j) + "/share_" + strconv.Itoa(i) + ".share")
							if err != nil {
								panic(err)
							}
							defer file.Close()

							file.Write(shares[byte(i)])
						}
					}
				} else {
					file, err := os.Create(dir + "/" + strconv.Itoa(i) + "/share_" + strconv.Itoa(i) + ".share")
					if err != nil {
						panic(err)
					}
					defer file.Close()

					file.Write(shares[byte(i)])
				}
			}

			// 一時停止
			fmt.Println("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			shares = map[byte][]byte{}

			// ファイル読み込み
			for i := 1; i <= n; i++ {
				file, err := os.Open(dir + "/" + strconv.Itoa(i) + "/share_" + strconv.Itoa(i) + ".share")
				if err != nil {
					log.Println("no such file or directory :" + dir + "/" + strconv.Itoa(i) + "/share_" + strconv.Itoa(i) + ".share")
					continue
				}
				defer file.Close()

				shares[byte(i)], err = ioutil.ReadAll(file)
				if err != nil {
					panic(err)
				}
			}

			log.Println(shares)

			// 分散シェアの復号化
			plain_shares := map[byte][]byte{}
			for i := 1; i <= n; i++ {
				plain_shares[byte(i)], err = ecies.Decrypt(keys[byte(i)], shares[byte(i)])
				if err != nil {
					log.Println("failed to decrypt ciphertext: ", err, "share_" + "/" + strconv.Itoa(i) + "/share_" + strconv.Itoa(i) + ".share")
					delete(plain_shares, byte(i))
				}
				log.Println("ciphertext decrypted: ", plain_shares[byte(i)])
			}

			decrypt_text := gr.Decrypt(plain_shares, block)
			log.Println(string(decrypt_text))

		default:
			log.Fatal("Error: No mode selected")
	}
}
