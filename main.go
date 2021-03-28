package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func convertFileToBytes(file string) []byte {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func convertBytesToFile(name string, bytesArr []byte, perm int) {
	err := ioutil.WriteFile(name, bytesArr, os.FileMode(perm))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func sendPacketTest(test []byte, size int) {
	i := 0
	buff := make([]byte, size)
	var str string = ""

	for {
		if i+size > len(test) {
			buff = test[i:]
			str = str + string(buff)
			fmt.Println("-----------------------")
			fmt.Println(string(test))
			fmt.Println(str)
			fmt.Println("-----------------------")
			fmt.Println("END")
			return
		}
		buff = test[i : i+size]
		str = str + string(buff)
		fmt.Println(string(buff))

		i = i + size
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	name := os.Args[1]

	listener, err := net.ListenPacket("udp4", name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("-------------------------")
	fmt.Printf("Listening on : %s\n", name)

	//s := "Hello Client ofjsdogsog sdyufsdug dfhsdgfifhieufhiw sdhfjjjjjjj dddd"
	str := ""
	//fileByte := []byte(s)
	//sendPacketTest(fileByte, 7)

	//fileByte := convertFileToBytes("Alpacas.jpeg")
	//fileByte := convertFileToBytes("test.jpg")
	fileByte := convertFileToBytes("test2.pdf")

	i := 0
	size := 1000

	buff := make([]byte, 1000)
	buff2 := make([]byte, size)
	for {
		//buff := make([]byte, 1000)

		n, conn, err := listener.ReadFrom(buff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
			continue
		}

		fmt.Println("-------------------------")

		log.Printf("Receive : %s\n", string(buff[:n]))
		if string(buff[:n]) == "OK" {
			if i+size > len(fileByte) {
				buff2 = fileByte[i:]
				str = string(buff2)
				listener.WriteTo([]byte(str), conn)
				log.Println("Final packet sent")
				//log.Printf("Send : %s\n", str)
				break
			}
			buff2 = fileByte[i : i+size]
			str = string(buff2)
			listener.WriteTo([]byte(str), conn)
			log.Println("Packet sent")
			//log.Printf("Send : %s\n", str)
		}

		fmt.Println("-------------------------")
		//handleClient(conn)
		i = i + size
	}
}
