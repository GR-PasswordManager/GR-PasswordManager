package gr

import (
	"fmt"
	"log"
	"regexp"

	"go.bug.st/serial"
)

const (
	// Raspberry PiのVIDとPID
	dongleVID string = "0525"
	donglePID string = "A4A7"
)

func Terminal(){
	// ポートをすべてスキャンし、指定されたVIDとPIDを持つポートを返す
	ports, err := getSerialPorts(dongleVID, donglePID)
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

	fmt.Printf("%sポートを指定して開きます\n", ports[0].Name)

	// シリアルポートを開く
	port, err := serial.Open(ports[0].Name, mode)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`\[.+?\]`)
	str := ""

	fmt.Printf("START:%q", ports[0].Name)

	for re.FindString(str) != "[dongle]" {
		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, "[who]\n")
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}

		// 受信したデータの出力
		fmt.Printf("T_Received data: '%q'EOF\n", str)
	}

	// SAVE
	str = ""
	for !re.MatchString(str) {
		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, "[save]")
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}
	}

	str = ""
	for !re.MatchString(str) {
		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, "[1]")
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}
	}

	str = ""
	for !re.MatchString(str) {
		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, "[share_abcdef]")
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("SAVE:%q", str)

	// シリアル通信でデータを送信する
	_, err = sendSerialData(port, "[save_complete]")
	if err != nil {
		log.Fatal(err)
	}

	// PICK

	str = ""
	for !re.MatchString(str) {
		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, "[pick]")
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}
	}

	str = ""
	for !re.MatchString(str) {
		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, "[1]")
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("PICK:%q", str)

	// シリアル通信でデータを送信する
	_, err = sendSerialData(port, "[pick_complete]")
	if err != nil {
		log.Fatal(err)
	}

	for re.FindString(str) != "[quit]" {
		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, "[quit]")
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}

		// 受信したデータの出力
		fmt.Printf("T_Received data: '%q'EOF\n", str)
	}

	// シリアルポートを閉じる
	err = port.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}
