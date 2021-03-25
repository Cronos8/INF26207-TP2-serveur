package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

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

	s := "Hello Client"

	for {
		buff := make([]byte, 1024)
		n, conn, err := listener.ReadFrom(buff)
		fmt.Println("-------------------------")
		log.Printf("Receive : %s\n", string(buff[:n]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
			continue
		}

		//buff[2] |= 0x80
		listener.WriteTo([]byte(s), conn)
		log.Printf("Send : %s\n", s)
		fmt.Println("-------------------------")
		//handleClient(conn)
	}
}
