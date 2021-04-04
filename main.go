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

func getByteSignature(packet []byte) [20]byte {
	return sha1.Sum(packet)
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

	fmt.Println("************************************************")
	fmt.Println()

	fmt.Printf("[Packet N° : %v]\n", nbPacket)
	fmt.Printf("Signature : %x\n", getByteSignature(packet))
	fmt.Printf("Corp du Packet - Signature : %x\n", getByteSignature(bodyPacket[8:]))

	fmt.Println()
	fmt.Println("************************************************")
}

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

	fileByte := convertFileToBytes("alpaga.jpeg")
	i := 0
	size := 1000
	buffRead := make([]byte, 1000)
	buffWrite := make([]byte, 1024)
	lastbuff := make([]byte, 1024)
	terminated := false

	var fiability float32 = 0.95
	var lastPacketConn net.Addr
	var nbPacket uint64 = 0

	if newClientConnexion(listener) == 0 {
		listener.SetReadDeadline(time.Now().Add(3 * time.Second))
		for terminated == false {
			n, conn, err := listener.ReadFrom(buffRead)
			switch e := isTimeOutError(err); e {
			case 1:
				_, err2 := listener.WriteTo(lastbuff, lastPacketConn)
				if err2 != nil {
					fmt.Println(err2)
				}
				fmt.Println("-------------------------")
				log.Println("RENVOIE PACKET")
				printPacket(lastbuff)
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

				if sendPaquetWithFiability(fiability) == true {
					if (string(buffRead[:n]) == "PACKAGE RECEIVE") || (string(buffRead[:n]) == "READY TO RECEIVE") {
						if i+size > len(fileByte) {
							buffWrite = encapPacket(nbPacket, fileByte[i:])
							listener.WriteTo(buffWrite, conn)
							fmt.Println("-------------------------")
							log.Println("DERNIER PACKET")
							printPacket(buffWrite)
							fmt.Println("-------------------------")
							fmt.Println()

							listener.WriteTo([]byte("END"), conn)
							getFileByteSignature(fileByte)
							terminated = true
							break
						}
						buffWrite = encapPacket(nbPacket, fileByte[i:i+size])
						listener.WriteTo(buffWrite, conn)
						fmt.Println("-------------------------")
						log.Println("PACKET SEND")
						printPacket(buffWrite)
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
