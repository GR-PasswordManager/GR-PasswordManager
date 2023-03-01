package gr

import (
	"fmt"
	"log"

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

	str := ""

	fmt.Printf("START:%q", ports[0].Name)

	checkSendSerialData(port, "who")
	str = checkReceiveSerialData(port)
	if str != "dongle" {
		log.Fatal("Dongle not found")
	}
	log.Println("dongle found")

	// SAVE
	checkSendSerialData(port, "save")
	checkSendSerialData(port, "1")
	checkSendSerialData(port, "share_abcdef")

	// PICK
	checkSendSerialData(port, "pick")
	checkSendSerialData(port, "1")
	str = checkReceiveSerialData(port)
	log.Printf("PICK:%q", str)

	// QUIT
	checkSendSerialData(port, "quit")
	// シリアルポートを閉じる
	err = port.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}
