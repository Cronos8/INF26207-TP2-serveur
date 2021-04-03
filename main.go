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

	//fileByte := convertFileToBytes("Alpacas.jpeg")
	fileByte := convertFileToBytes("test.jpg")
	//fileByte := convertFileToBytes("test2.pdf")

	i := 0
	size := 1000

	buff := make([]byte, 1000)
	//buff2 := make([]byte, size)
	buff2 := make([]byte, 1024)
	var lastPacketConn net.Addr
	lastbuff := make([]byte, 1024)
	terminated := false
	var nbPacket uint64 = 0

	if newClientConnexion(listener) == 0 {
		listener.SetReadDeadline(time.Now().Add(3 * time.Second))
		for terminated == false {
			//buff := make([]byte, 1000)
			n, conn, err := listener.ReadFrom(buff)
			//fmt.Println(n)
			switch e := isTimeOutError(err); e {
			case 1:
				_, err3 := listener.WriteTo(lastbuff, lastPacketConn)
				if err3 != nil {
					//fmt.Println(conn.Network())
					//fmt.Println(conn.String())
					fmt.Println("error1")
					fmt.Println(err3)
				}

				log.Println("RENVOIE PACKET")
				log.Println("Packet entier")
				getByteSignature(lastbuff)
				log.Println("numero Packet")
				getByteSignature(lastbuff[:8])
				log.Println("corp du Packet")
				getByteSignature(lastbuff[8:])
				log.Println("Packet re-sent")
				listener.SetReadDeadline(time.Now().Add(3 * time.Second))
				break
			case -1:
				fmt.Println("erro3")
				fmt.Println(e)
				continue
			case 0:
				fmt.Println("-------------------------")

				log.Printf("Receive : %s\n", string(buff[:n]))
				if sendPaquetWithFiability(0.95) == true {
					if (string(buff[:n]) == "PACKAGE RECEIVE") || (string(buff[:n]) == "READY TO RECEIVE") {
						if i+size > len(fileByte) {
							buff2 = encapPacket(nbPacket, fileByte[i:])
							//buff2 = fileByte[i:]
							listener.WriteTo(buff2, conn)
							log.Println("Final packet sent")
							log.Printf("Packet nb : %v\n", nbPacket)
							getByteSignature(buff2)
							log.Println("numero Packet")
							getByteSignature(buff2[:8])
							log.Println("corp du Packet")
							getByteSignature(buff2[8:])
							//log.Printf("Send : %s\n", str)
							listener.WriteTo([]byte("END"), conn)
							getFileByteSignature(fileByte)
							terminated = true
							break
						}
						buff2 = encapPacket(nbPacket, fileByte[i:i+size])
						//buff2 = fileByte[i : i+size]
						listener.WriteTo(buff2, conn)
						log.Println("Packet sent")
						log.Printf("Packet nb : %v\n", nbPacket)
						log.Println("Packet entier")
						getByteSignature(buff2)
						log.Println("numero Packet")
						getByteSignature(buff2[:8])
						log.Println("corp du Packet")
						getByteSignature(buff2[8:])

						nbPacket++
						//log.Printf("Send : %s\n", str)
					}

					fmt.Println("-------------------------")
					i = i + size
					lastbuff = buff2
					//copy(lastbuff, buff2)
					lastPacketConn = conn
					fmt.Println("************************")
					getByteSignature(lastbuff)
					getByteSignature(buff2)
					fmt.Println("************************")
					//lastbuff = buff2
					listener.SetReadDeadline(time.Now().Add(3 * time.Second))

				} else {
					log.Println("Fiability Error")
				}
			}
		}
	}
}
