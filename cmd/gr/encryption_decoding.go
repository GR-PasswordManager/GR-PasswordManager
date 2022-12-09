package gr

import (
	"log"

	"github.com/codahale/sss"
	ecies "github.com/ecies/go/v2"
)

func Encrypt(n int, k int, secret string, publickeys []*ecies.PublicKey) ([]*ecies.PrivateKey, [][]byte) {
	// 分散シェアの生成
	shares, err := sss.Split(byte(n), byte(k), []byte(secret))
	if err != nil {
		panic(err)
	}
	log.Println("shares: ", shares)

	// 分散シェアの暗号化
	var encrypted_share [][]byte = nil
	keys := make([]*ecies.PrivateKey, n)

	if publickeys == nil {
		log.Println("publickeys is nil")
		publickeys = make([]*ecies.PublicKey, n)
		for i := 0; i < n; i++ {
			keys[i], err = ecies.GenerateKey()
			if err != nil {
				panic(err)
			}
			log.Println("key[", i, "] pair has been generated")
			publickeys[i] = keys[i].PublicKey
		}
	}

	for i := 0; i < n; i++ {
		ciphertext, err := ecies.Encrypt(publickeys[i], shares[byte(i+1)])
		if err != nil {
			panic(err)
		}
		encrypted_share = append(encrypted_share, ciphertext)
	}

	return keys, encrypted_share
}
