package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Cronos8/INF26207-TP2-serveur/filebyte"
	"github.com/Cronos8/INF26207-TP2-serveur/packet"
	"github.com/Cronos8/INF26207-TP2-serveur/serverfunc"
)

func checkArguments(args []string) int {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
		return -1
	}
	return 0
}

func main() {
	if checkArguments(os.Args) != 0 {
		return
	}

	listener, err := net.ListenPacket("udp4", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s\n", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("-------------------------")
	fmt.Printf("Listening on : %s\n", os.Args[1])
	fmt.Println("-------------------------")
	fmt.Println()

	fByte := filebyte.ConvertFileToBytes("testfiles/alpaga.jpeg")

	i := 0
	size := 1000
	buffRead := make([]byte, 1000)
	buffWrite := make([]byte, 1024)
	lastbuff := make([]byte, 1024)
	terminated := false

	var fiability float32 = 0.95
	var lastPacketConn net.Addr
	var nbPacket uint64 = 0

	if serverfunc.NewClientConnexion(listener) == 0 {
		listener.SetReadDeadline(time.Now().Add(3 * time.Second))
		for terminated == false {
			n, conn, err := listener.ReadFrom(buffRead)
			switch e := serverfunc.IsTimeOutError(err); e {
			case 1:
				_, err2 := listener.WriteTo(lastbuff, lastPacketConn)
				if err2 != nil {
					fmt.Println(err2)
				}
				fmt.Println("-------------------------")
				log.Println("RENVOIE PACKET")
				packet.PrintPacket(lastbuff)
				fmt.Println("-------------------------")
				fmt.Println()
				listener.SetReadDeadline(time.Now().Add(3 * time.Second))
				break

			case -1:
				fmt.Println(e)
				break

			case 0:

				fmt.Println("-------------------------")
				log.Printf("Receive : %s\n", string(buffRead[:n]))
				fmt.Println("-------------------------")
				fmt.Println()

				if serverfunc.SendPaquetWithFiability(fiability) == true {
					if (string(buffRead[:n]) == "PACKAGE RECEIVE") || (string(buffRead[:n]) == "READY TO RECEIVE") {
						if i+size > len(fByte) {
							buffWrite = packet.EncapPacket(nbPacket, fByte[i:])
							listener.WriteTo(buffWrite, conn)
							fmt.Println("-------------------------")
							log.Println("DERNIER PACKET")
							packet.PrintPacket(buffWrite)
							fmt.Println("-------------------------")
							fmt.Println()

							listener.WriteTo([]byte("END"), conn)
							filebyte.GetFileByteSignature(fByte)
							terminated = true
							break
						}
						buffWrite = packet.EncapPacket(nbPacket, fByte[i:i+size])
						listener.WriteTo(buffWrite, conn)
						fmt.Println("-------------------------")
						log.Println("PACKET SEND")
						packet.PrintPacket(buffWrite)
						fmt.Println("-------------------------")
						fmt.Println()
						nbPacket++
					}

					i = i + size
					lastbuff = buffWrite
					lastPacketConn = conn
					listener.SetReadDeadline(time.Now().Add(3 * time.Second))

				} else {
					fmt.Println("!!!!!!!!!!!!!!!!!!")
					log.Println("Fiability Error")
					fmt.Println("!!!!!!!!!!!!!!!!!!")
					fmt.Println()
				}
			}
		}
	}
}
