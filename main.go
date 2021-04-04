package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

// uint64 -> 8 octets
func encapPacket(nbPacket uint64, packet []byte) []byte {
	buffnbpacket := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffnbpacket, nbPacket)

	buffdest := append(buffnbpacket, packet...)

	return buffdest
}

func sendPaquetWithFiability(fiability float32) bool {
	if rand.Float32() <= fiability {
		return true
	}
	return false
}

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

func getByteSignature(packet []byte) {
	log.Printf("Packet signature : %x\n", sha1.Sum(packet))
}

func getFileSignature(file string) {
	data := convertFileToBytes(file)
	fmt.Printf("File signature : %x\n", sha1.Sum(data))
}

func getFileByteSignature(fileByte []byte) {
	log.Printf("File signature : %x\n", sha1.Sum(fileByte))
}

func newClientConnexion(listener net.PacketConn) int {
	buff := make([]byte, 1000)

	listener.SetReadDeadline(time.Now().Add(5 * time.Minute))
	for {
		n, co, err := listener.ReadFrom(buff)
		if e, ok := err.(net.Error); ok && e.Timeout() {
			fmt.Fprintf(os.Stderr, "Timeout Error : %s\n", err.Error())
			return -1
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal Error : %s\n", err.Error())
			continue
		}

		if string(buff[:n]) == "Client - CONNEXION OK" {
			listener.WriteTo([]byte("Serveur - CONNEXION OK"), co)
			break
		}
	}
	log.Println("Connexion established")
	return 0
}

func isTimeOutError(err error) int {
	if e, ok := err.(net.Error); ok && e.Timeout() {
		fmt.Fprintf(os.Stderr, "Timeout Error : %s\n", err.Error())
		return 1
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error : %s\n", err.Error())
		return -1
	}
	return 0
}

func decapPacket(packet []byte) (uint64, []byte) {
	buffnbpacket := packet[:8]
	buffbody := packet[8:]
	nbPacket := binary.LittleEndian.Uint64(buffnbpacket)

	return nbPacket, buffbody
}

func printPacket(packet []byte) {
	nbPacket, bodyPacket := decapPacket(packet)
	log.Printf("Packet nb : %v\n", nbPacket)

	log.Println("Packet entier")
	getByteSignature(bodyPacket)

	log.Println("numero Packet")
	getByteSignature(bodyPacket[:8])

	log.Println("corp du Packet")
	getByteSignature(bodyPacket[8:])
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	name := os.Args[1]

	listener, err := net.ListenPacket("udp4", name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s\n", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("-------------------------")
	fmt.Printf("Listening on : %s\n", name)

	fileByte := convertFileToBytes("alpaga.jpeg")
	i := 0
	size := 1000
	buff := make([]byte, 1000)
	buff2 := make([]byte, 1024)
	var lastPacketConn net.Addr
	lastbuff := make([]byte, 1024)
	terminated := false
	var nbPacket uint64 = 0

	if newClientConnexion(listener) == 0 {
		listener.SetReadDeadline(time.Now().Add(3 * time.Second))
		for terminated == false {
			n, conn, err := listener.ReadFrom(buff)
			switch e := isTimeOutError(err); e {
			case 1:
				_, err2 := listener.WriteTo(lastbuff, lastPacketConn)
				if err2 != nil {
					fmt.Println(err2)
				}
				log.Println("RENVOIE PACKET")
				printPacket(lastbuff)
				log.Println("Packet re-sent")
				listener.SetReadDeadline(time.Now().Add(3 * time.Second))
				break

			case -1:
				fmt.Println(e)
				break

			case 0:
				fmt.Println("-------------------------")

				log.Printf("Receive : %s\n", string(buff[:n]))
				if sendPaquetWithFiability(0.95) == true {
					if (string(buff[:n]) == "PACKAGE RECEIVE") || (string(buff[:n]) == "READY TO RECEIVE") {
						if i+size > len(fileByte) {
							buff2 = encapPacket(nbPacket, fileByte[i:])
							listener.WriteTo(buff2, conn)
							log.Println("Final packet sent")
							printPacket(buff2)

							listener.WriteTo([]byte("END"), conn)
							getFileByteSignature(fileByte)
							terminated = true
							break
						}
						buff2 = encapPacket(nbPacket, fileByte[i:i+size])
						listener.WriteTo(buff2, conn)
						log.Println("Packet sent")
						printPacket(buff2)
						nbPacket++
					}

					fmt.Println("-------------------------")
					i = i + size
					lastbuff = buff2
					lastPacketConn = conn
					listener.SetReadDeadline(time.Now().Add(3 * time.Second))

				} else {
					log.Println("Fiability Error")
				}
			}
		}
	}
}
