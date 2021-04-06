package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Cronos8/INF26207-TP2-serveur/filebyte"
	"github.com/Cronos8/INF26207-TP2-serveur/packet"
	"github.com/Cronos8/INF26207-TP2-serveur/serverfunc"
)

func checkArguments(args []string) int {
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr filepath\n", os.Args[0])
		os.Exit(1)
		return -1
	}
	return 0
}

func main() {

	hpacket := packet.HeaderPacket{nil, 0, 0}

	if checkArguments(os.Args) != 0 {
		return
	}

	udpaddr, err := net.ResolveUDPAddr("udp4", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid IPv4 address...\n")
		os.Exit(1)
	}

	hpacket.HeaderIp = udpaddr.IP
	hpacket.HeaderPort = int32(udpaddr.Port)

	// "testfiles/alpaga.jpeg"
	fByte := filebyte.ConvertFileToBytes(os.Args[2])
	if fByte == nil {
		os.Exit(1)
	}
	listener, err := net.ListenPacket("udp4", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(listener.LocalAddr().String())
	defer listener.Close()

	fmt.Println("------------------------------------------------")
	fmt.Printf("Listening on : %s\n", os.Args[1])
	fmt.Println("------------------------------------------------")
	fmt.Println()

	i := 0
	//size := 1000
	size := 996
	buffRead := make([]byte, 1000)
	buffWrite := make([]byte, 1024)
	lastbuff := make([]byte, 1024)
	terminated := false

	var fiability float32 = 0.95
	var nbPacket uint64 = 0

	conn := serverfunc.NewClientConnexion(listener)
	if conn != nil {
		listener.SetReadDeadline(time.Now().Add(3 * time.Second))
		for terminated == false {
			n, _, err := listener.ReadFrom(buffRead)
			//fmt.Println(conn.String())
			switch e := serverfunc.IsTimeOutError(err); e {
			case 1:
				_, err2 := listener.WriteTo(lastbuff, conn)
				if err2 != nil {
					fmt.Println(err2)
				}
				packet.PrintMessage("PACKET RE-SEND", packet.CyanColor, conn.String())
				packet.PrintPacket(lastbuff)
				fmt.Println()
				listener.SetReadDeadline(time.Now().Add(3 * time.Second))
				break

			case -1:
				fmt.Println(e)
				break

			case 0:

				packet.PrintMessage("MESSAGE RECEIVE", packet.GreenColor, conn.String())
				fmt.Printf("Content : %s\n", string(buffRead[:n]))

				if serverfunc.SendPaquetWithFiability(fiability) == true {
					if (string(buffRead[:n]) == "PACKAGE RECEIVE") || (string(buffRead[:n]) == "READY TO RECEIVE") {
						if i+size > len(fByte) {
							buffWrite = packet.EncapPacket(hpacket, fByte[i:])
							listener.WriteTo(buffWrite, conn)
							packet.PrintMessage("LAST PACKET", packet.GreenColor, conn.String())
							packet.PrintPacket(buffWrite)
							fmt.Println()

							listener.WriteTo([]byte("END"), conn)
							filebyte.GetFileByteSignature(fByte)
							terminated = true
							break
						}
						buffWrite = packet.EncapPacket(hpacket, fByte[i:i+size])

						listener.WriteTo(buffWrite, conn)
						packet.PrintMessage("PACKET SEND", packet.BlueColor, conn.String())
						packet.PrintPacket(buffWrite)
						fmt.Println()
						nbPacket++
						hpacket.HeaderNbPacket = nbPacket
					}

					i = i + size
					lastbuff = buffWrite
					listener.SetReadDeadline(time.Now().Add(3 * time.Second))

				} else {
					packet.PrintMessage("FIABILITY ERROR", packet.RedColor, conn.String())
				}
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Server timeout ...")
		os.Exit(1)
	}
}
