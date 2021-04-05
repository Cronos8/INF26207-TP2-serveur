package packet

import (
	"encoding/binary"
	"fmt"
	"log"
	"github.com/Cronos8/INF26207-TP2-serveur/filebyte"
)

// ColorPrint colors print
type ColorPrint string 

const (
	BlueColor ColorPrint = "\033[34m"
	RedColor ColorPrint = "\033[31m"
	GreenColor ColorPrint = "\033[32m"
	ResetColor ColorPrint = "\033[0m"
	YellowColor ColorPrint = "\033[33m"
	CyanColor ColorPrint = "\033[36m"
	
	/*
    colorPurple := "\033[35m"
	colorWhite := "\033[37m"
	*/
)

// EncapPacket packet encapsulation uint64 -> 8 octets
func EncapPacket(nbPacket uint64, packet []byte) []byte {
	buffnbpacket := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffnbpacket, nbPacket)

	buffdest := append(buffnbpacket, packet...)

	return buffdest
}

// DecapPacket packet decapsulation
func DecapPacket(packet []byte) (uint64, []byte) {
	buffnbpacket := packet[:8]
	buffbody := packet[8:]
	nbPacket := binary.LittleEndian.Uint64(buffnbpacket)

	return nbPacket, buffbody
}

// PrintMessage print a message 
func PrintMessage(message string, color ColorPrint){
	fmt.Println(string(color))
	fmt.Println("------------------------------------------------")
	log.Println(message)
	fmt.Println("------------------------------------------------")
	fmt.Println(string(ResetColor))
}

// PrintPacket print a packet
func PrintPacket(p []byte) {

	nbPacket, bodyPacket := DecapPacket(p)

	fmt.Printf("\t************************************************\n")
	fmt.Println()

	fmt.Printf("\t[Packet N° : %v]\n", nbPacket)
	fmt.Printf("\tSignature : %x\n", filebyte.GetByteSignature(p))
	fmt.Printf("\tCorp du Packet - Signature : %x\n", filebyte.GetByteSignature(bodyPacket[8:]))

	fmt.Println()
	fmt.Printf("\t************************************************\n")
}
