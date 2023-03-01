package gr

import (
	"errors"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

// ポートをすべてスキャンし、指定されたVIDとPIDを持つポートを返す
// VIDとPIDが空の場合はすべてのポートを返す
func getSerialPorts(VID string, PID string) ([]*enumerator.PortDetails, error) {

	// ポートをすべてスキャンする
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}

	// 接続されたものがなかった場合
	if len(ports) == 0 {
		log.Println("No serial ports found!")
		err = errors.New("No serial ports found!")
		return nil, err
	}

	// 指定されたVIDとPIDを持つポートを探す
	var serialPorts []*enumerator.PortDetails
	for _, port := range ports {
		log.Printf("Found port: %s(VID:%s, PID:%s)\n", port.Name, port.VID, port.PID)
		// VIDとPIDが一致するかどうか
		if port.VID == VID && port.PID == PID || VID == "" && PID == ""{
			serialPorts = append(serialPorts, port)
			log.Printf("Found!\n")
			log.Printf("   Name       %s\n", port.Name)
			log.Printf("    ID        %s:%s\n", port.VID, port.PID)
			log.Printf("   USB serial %s\n", port.SerialNumber)
		}
	}

	// デバイスは接続されているが、指定のものが接続されていなかった場合
	if len(serialPorts) == 0 {
		log.Println("No serial ports found!")
		err = errors.New("No serial ports found!")
		return nil, err
	}

	return serialPorts, err
}

func sendSerialData(port serial.Port, data string) (int, error) {
	// 送信するデータの出力
	log.Printf("Sending data: %q\n", data)

	// シリアル通信でデータを送信する
	n, err := port.Write([]byte(data + "\n\r"))
	return n, err
}

func receiveSerialData(port serial.Port) (string, error) {
	// 受信したデータの全体を格納する変数
	data := ""

	port.SetReadTimeout(1000 * time.Millisecond)

	// 受信するデータのバッファ先を作成する
	buff := make([]byte, 8)
	for {
		// 作成したバッファ分のデータを受信する
		n, err := port.Read(buff)
		if err != nil {
			break
		}
		// もし、データがなければループを抜ける
		if n == 0 {
			log.Println("取得したデータが空です")
			break
		}

		// 受信したデータを格納する
		data += string(buff[:n])

		// 受信したデータに"\n"が含まれていたらループを抜ける
		if strings.Contains(string(buff[:n]), "\n") {
			log.Println("改行を検出")
			break
		}
	}

	// 受信したデータの出力
	log.Printf("Received data: %q\n", data)
	return data, nil
}

func checkSendSerialData(port serial.Port, data string) {
	str := ""

	for data != ("c_" + str) {
		// シリアル通信でデータを送信する
		_, err := sendSerialData(port, data)
		if err != nil {
			log.Fatal(err)
		}

		// シリアル通信でデータを受信する
		str, err = receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("str: %s\n", str)
		log.Printf("data: %s\n", data)
	}

	// シリアル通信でデータを送信する
	_, err := sendSerialData(port, "OK")
	if err != nil {
		log.Fatal(err)
	}
}

func checkReceiveSerialData(port serial.Port) (string) {
	str := ""

	for {
		// シリアル通信でデータを受信する
		str, err := receiveSerialData(port)
		if err != nil {
			log.Fatal(err)
		}

		if str == "OK" {
			break
		}

		// シリアル通信でデータを送信する
		_, err = sendSerialData(port, ("c_" + str))
		if err != nil {
			log.Fatal(err)
		}
	}
	return str
}
