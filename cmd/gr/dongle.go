package gr

import (
	"fmt"
	"log"
	"regexp"

	"go.bug.st/serial"
)

func Dongle(){
	// ポートをすべてスキャンし、指定されたVIDとPIDを持つポートを返す
	ports, err := getSerialPorts("", "")
	if err != nil {
		log.Fatal(err)
	}

	// シリアルポートの設定を行う
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	// シリアルポートを開く
	port, err := serial.Open(ports[0].Name, mode)
	if err != nil {
		log.Fatal(err)
	}

	str := "" // 何かしらの文字列を入れておく
	re := regexp.MustCompile(`\[.+?\]`)
	for {
		for !re.MatchString(str) {
			// シリアル通信でデータを受信する
			str, err = receiveSerialData(port)
			if err != nil {
				log.Fatal(err)
			}
		}

		switch re.FindString(str) {
			case "[who]":
				// シリアル通信でデータを送信する
				_, err = sendSerialData(port, "[dongle]")
				if err != nil {
					log.Fatal(err)
				}

			case "[test]":
				// シリアル通信でデータを送信する
				_, err = sendSerialData(port, "[test_d]")
				if err != nil {
					log.Fatal(err)
				}

			case "[quit]":
				// シリアル通信でデータを送信する
				_, err = sendSerialData(port, "[quit]")
				if err != nil {
					log.Fatal(err)
				}
				break

			default:
				// 受信したデータの出力
				fmt.Printf("Received data: %s", re.FindAllString(str, -1))
		}
	}

	// シリアルポートを閉じる
	err = port.Close()
	if err != nil {
		log.Fatal(err)
	}
}
