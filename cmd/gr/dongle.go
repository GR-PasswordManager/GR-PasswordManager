package gr

import (
	"fmt"
	"log"
	"regexp"

	"go.bug.st/serial"
)

func Dongle(){
	// シリアルポートの設定を行う
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port_name := "/dev/ttyGS0"

	fmt.Printf("%sポートを指定して開きます\n", port_name)

	// シリアルポートを開く
	port, err := serial.Open(port_name, mode)
	if err != nil {
		log.Fatal(err)
	}

	str := "" // 何かしらの文字列を入れておく
	re := regexp.MustCompile(`\[.+?\]`)

	fmt.Printf("START:%q", port_name)

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
				fmt.Printf("D_Received data: '%q'EOF\n", re.FindAllString(str, -1))
		}
	}

	// シリアルポートを閉じる
	err = port.Close()
	if err != nil {
		log.Fatal(err)
	}
}
