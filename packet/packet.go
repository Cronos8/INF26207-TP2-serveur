package packet

import (
	"encoding/binary"
	"fmt"
	"github.com/Cronos8/INF26207-TP2-serveur/filebyte"
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

// PrintPacket print a packet
func PrintPacket(p []byte) {

	nbPacket, bodyPacket := DecapPacket(p)

	fmt.Println("************************************************")
	fmt.Println()

	fmt.Printf("[Packet N° : %v]\n", nbPacket)
	fmt.Printf("Signature : %x\n", filebyte.GetByteSignature(p))
	fmt.Printf("Corp du Packet - Signature : %x\n", filebyte.GetByteSignature(bodyPacket[8:]))

	fmt.Println()
	fmt.Println("************************************************")
}
