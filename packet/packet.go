package packet

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/Cronos8/INF26207-TP2-serveur/filebyte"
)

// ColorPrint colors print
type ColorPrint string

const (
	BlueColor   ColorPrint = "\033[34m"
	RedColor    ColorPrint = "\033[31m"
	GreenColor  ColorPrint = "\033[32m"
	ResetColor  ColorPrint = "\033[0m"
	YellowColor ColorPrint = "\033[33m"
	CyanColor   ColorPrint = "\033[36m"

	/*
		    colorPurple := "\033[35m"
			colorWhite := "\033[37m"
	*/
)

// HeaderPacket header of packet
type HeaderPacket struct {
	HeaderIp       net.IP // 16 byte - 128 ctets
	HeaderPort     int32  // 4 byte -> 32 octets
	HeaderNbPacket uint64 // 8 byte -> 64 octets
	// // 28 bytes au total
}

// EncapPacket2 packet encapsulation uint64 -> 8 octets
func EncapPacket2(nbPacket uint64, packet []byte) []byte {

	buffnbpacket := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffnbpacket, nbPacket)

	buffdest := append(buffnbpacket, packet...)

	return buffdest
}

// DecapPacket2 packet decapsulation
func DecapPacket2(packet []byte) (uint64, []byte) {
	buffnbpacket := packet[:8]
	buffbody := packet[8:]
	nbPacket := binary.LittleEndian.Uint64(buffnbpacket)

	return nbPacket, buffbody
}

// EncapPacket encapsul packet
func EncapPacket(hpacket HeaderPacket, packet []byte) []byte {

	buffnbpacket := make([]byte, 8)
	buffip := make([]byte, 16)
	buffport := make([]byte, 4)

	buffip = []byte(hpacket.HeaderIp) // buffip len = 16 et non 4
	fmt.Println(len(buffip))

	binary.LittleEndian.PutUint32(buffport, uint32(hpacket.HeaderPort))
	binary.LittleEndian.PutUint64(buffnbpacket, hpacket.HeaderNbPacket)

	bufftmp := []byte{}
	bufftmp = append(bufftmp, buffnbpacket...)
	bufftmp = append(bufftmp, buffip...)
	bufftmp = append(bufftmp, buffport...)
	bufftmp = append(bufftmp, packet...)

	return bufftmp
}

// DecapPacket packet decapsulation
func DecapPacket(packet []byte) (HeaderPacket, []byte) {

	buffnbpacket := packet[:8]
	nbPacket := binary.LittleEndian.Uint64(buffnbpacket)
	buffipacket := packet[8:24]
	buffportpacket := packet[24:28]
	nbPort := binary.LittleEndian.Uint32(buffportpacket)

	hpacket := HeaderPacket{
		net.IP(buffipacket),
		int32(nbPort),
		nbPacket,
	}
	buffbody := packet[28:]

	return hpacket, buffbody
}

// PrintMessage print a message
func PrintMessage(message string, color ColorPrint) {
	fmt.Println(string(color))
	fmt.Println("-----------------------------------------")
	log.Println(message)
	fmt.Println("---------------------------------------------")
	fmt.Println(string(ResetColor))
}

// PrintPacket print a packet
func PrintPacket(p []byte) {

	hpacket, bodyPacket := DecapPacket(p)

	fmt.Printf("\t************************************************\n")
	fmt.Println()

	fmt.Printf("\t[Packet N° : %v]\n", hpacket.HeaderNbPacket)
	fmt.Printf("\tSignature : %x\n", filebyte.GetByteSignature(p))
	fmt.Printf("\tCorp du Packet - Signature : %x\n", filebyte.GetByteSignature(bodyPacket))

	fmt.Println()
	fmt.Printf("\t************************************************\n")
}
