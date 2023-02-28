package gr

import (
	"errors"
	"fmt"
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
		fmt.Println("No serial ports found!")
		err = errors.New("No serial ports found!")
		return nil, err
	}

	// 指定されたVIDとPIDを持つポートを探す
	var serialPorts []*enumerator.PortDetails
	for _, port := range ports {
		fmt.Printf("Found port: %s(VID:%s, PID:%s)\n", port.Name, port.VID, port.PID)
		// VIDとPIDが一致するかどうか
		if port.VID == VID && port.PID == PID || VID == "" && PID == ""{
			serialPorts = append(serialPorts, port)
			fmt.Printf("Found!\n")
			fmt.Printf("   Name       %s\n", port.Name)
			fmt.Printf("    ID        %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
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
	// 送信するデータの出力
	fmt.Printf("Sending data: '%s'EOF\n", data)

	// シリアル通信でデータを送信する
	n, err := port.Write([]byte(data + "\n\r"))
	return n, err
}

func receiveSerialData(port serial.Port) (string, error) {
	// 受信したデータの全体を格納する変数
	data := ""

	port.SetReadTimeout(500 * time.Millisecond)

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
			break
		}

		// 受信したデータを格納する
		data += string(buff[:n])

		// 受信したデータに"\n"が含まれていたらループを抜ける
		if strings.Contains(string(buff[:n]), "\n") {
			break
		}
	}

	// 受信したデータの出力
	fmt.Printf("Received data: '%s'EOF\n", data)
	return data, nil
}
