package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const (
	// ArduinoのVIDとPID
	VID string = "2341"
	PID string = "0001"
)

// ポートをすべてスキャンし、指定されたVIDとPIDを持つポートを返す
func getSerialPorts() ([]*enumerator.PortDetails, error) {

	// ポートをすべてスキャンする
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}

	// 接続されたものがなかった場合
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		err = errors.New("No serial ports found!")
		return nil, err
	}

	// 指定されたVIDとPIDを持つポートを探す
	var serialPorts []*enumerator.PortDetails
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		// 接続されているポートがUSBポートであるかどうか
		if port.IsUSB {
			// VIDとPIDが一致するかどうか
			if port.VID == VID && port.PID == PID {
				serialPorts = append(serialPorts, port)
				fmt.Printf("Found Arduino!\n")
				fmt.Printf("   Name       %s\n", port.Name)
				fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
				fmt.Printf("   USB serial %s\n", port.SerialNumber)
			}
		}
	}

	// デバイスは接続されているが、指定のものが接続されていなかった場合
	if len(serialPorts) == 0 {
		fmt.Println("No serial ports found!")
		err = errors.New("No serial ports found!")
		return nil, err
	}

	return serialPorts, err
}

func sendSerialData(port serial.Port, data string) (int, error) {
	// シリアル通信でデータを送信する
	n, err := port.Write([]byte(data + "\n\r"))
	return n, err
}

func receiveSerialData(port serial.Port) (string, error) {
	// 受信したデータの全体を格納する変数
	data := ""

	// 受信するデータのバッファ先を作成する
	buff := make([]byte, 4)
	for {
		// 作成したバッファ分のデータを受信する
		n, err := port.Read(buff)
		if err != nil {
			break
		}
		// もし、データがなければループを抜ける
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		// 受信したデータを格納する
		data += string(buff[:n])

		// 受信したデータに"\n"が含まれていたらループを抜ける
		if strings.Contains(string(buff[:n]), "\n") {
			break
		}
	}
	return data, nil
}

func main(){
	// ポートをすべてスキャンし、指定されたVIDとPIDを持つポートを返す
	ports, err := getSerialPorts()
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

	for {
		// シリアル通信でデータを送信する
		n, err := sendSerialData(port, "Hello World!")
		if err != nil {
			log.Fatal(err)
		}

		// 送信したバイト数を表示する
		fmt.Printf("Sent %v bytes\n", n)

		// シリアル通信でデータを受信する
		str, err := receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}

		// 受信したデータの出力
		fmt.Printf("Received data: %s", str)
	}
}
