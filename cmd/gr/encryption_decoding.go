package gr

import (
	"log"

	"github.com/codahale/sss"
	ecies "github.com/ecies/go/v2"
)

func Encrypt(n int, k int, secret string, publickeys map[byte]*ecies.PublicKey) (map[byte]*ecies.PrivateKey, map[byte][]byte) {
	// 分散シェアの生成
	shares, err := sss.Split(byte(n), byte(k), []byte(secret))
	if err != nil {
		panic(err)
	}
	log.Println("shares: ", shares)

	// 分散シェアの暗号化
	encrypted_share := map[byte][]byte{}
	keys := map[byte]*ecies.PrivateKey{}

	if publickeys == nil {
		log.Println("publickeys is nil")
		publickeys = map[byte]*ecies.PublicKey{}
		for i := 1; i <= n; i++ {
			keys[byte(i)], err = ecies.GenerateKey()
			if err != nil {
				panic(err)
			}
			log.Println("key[", i, "] pair has been generated")
			publickeys[byte(i)] = keys[byte(i)].PublicKey
		}
	}

	for i := 1; i <= n; i++ {
		ciphertext, err := ecies.Encrypt(publickeys[byte(i)], shares[byte(i)])
		if err != nil {
			panic(err)
		}
		encrypted_share[byte(i)] = ciphertext
	}

	return keys, encrypted_share
}
