package gr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"

	"github.com/codahale/sss"
	ecies "github.com/ecies/go/v2"
)

func Encrypt(k int, n int, secret []byte, block cipher.Block, publickeys map[byte]*ecies.PublicKey) (map[byte][]byte) {

	// chipherがない場合はエラー
	if block == nil {
		panic("chipher is nil")
	}

	// Create IV
	cipherText := make([]byte, aes.BlockSize+len(secret))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("err: %s\n", err)
	}

	// Encrypt
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], secret)
	log.Printf("Cipher text: %x \n", cipherText)

	// 分散シェアの生成
	shares, err := sss.Split(byte(n), byte(k), []byte(cipherText))
	if err != nil {
		panic(err)
	}
	log.Println("shares: ", shares)

	// 暗号化したシェアを格納するための配列
	encrypted_share := map[byte][]byte{}

	// 公開鍵がない場合はエラー
	if publickeys == nil {
		panic("publickeys is nil")
	}

	// 公開鍵での暗号化
	for i := 1; i <= n; i++ {
		ciphertext, err := ecies.Encrypt(publickeys[byte(i)], shares[byte(i)])
		if err != nil {
			panic(err)
		}
		encrypted_share[byte(i)] = ciphertext
	}

	return encrypted_share
}

func Decrypt(shares map[byte][]byte, block cipher.Block) []byte {
	log.Println(shares)

	// 分散シェアの結合
	cipherText := sss.Combine(shares)
	log.Printf("Combine_cipherText: %x \n", cipherText)

	// Decrpt
	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])

	return decryptedText
}
