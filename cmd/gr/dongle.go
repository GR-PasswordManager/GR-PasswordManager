package gr

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	// ディレクトリ作成
	dir := "share"
	err = os.Mkdir(dir, 0777)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	fmt.Printf("START:%q", port_name)

	quit:
		for {
			for !re.MatchString(str) {
				// シリアル通信でデータを受信する
				str = checkReceiveSerialData(port)
			}

			switch re.FindString(str) {
				case "[who]":
					// シリアル通信でデータを送信する
					checkSendSerialData(port, "[dongle]")

				case "[save]":

					// シェア名の受信
					share_name := checkReceiveSerialData(port)
					share_data := checkReceiveSerialData(port)

					log.Printf("share_name:%s", share_name)

					// 受信したデータのファイルへの書き込み
					file, err := os.Create(dir + "/" + share_name + ".share")
					if err != nil {
						panic(err)
					}
					defer file.Close()

					file.Write([]byte(share_data))

					log.Println("save complete")

				case "[pick]":
					// シェア名の受信
					share_name := checkReceiveSerialData(port)

					// 受信した名前のシェアファイルを開く
					share := []byte{}
					file, err := os.Open(dir + "/" + share_name + ".share")
					if err != nil {
						log.Println("no such file or directory :" + dir + "/" + share_name + ".share")
						share = []byte("[no_share]")
					} else {
						// ファイルから取り出し
						share, err = ioutil.ReadAll(file)
						if err != nil {
							panic(err)
						}
					}
					defer file.Close()

					str = ""
					for !re.MatchString(str) {
						// シリアル通信でデータを送信する
						_, err = sendSerialData(port, string(share))
						if err != nil {
							log.Fatal(err)
						}

						// シリアル通信でデータを受信する
						str, err = receiveSerialData(port)
						if err != nil {
							log.Fatal(err)
						}
					}

					log.Println("pick complete")

				case "[quit]":
					break quit

				default:
					// 受信したデータの出力
					fmt.Printf("D_Received data: '%q'EOF\n", re.FindAllString(str, -1))
			}
			str = ""
		}

	// シリアルポートを閉じる
	err = port.Close()
	if err != nil {
		log.Fatal(err)
	}
}
