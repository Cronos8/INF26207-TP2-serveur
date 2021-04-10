package serverfilebyte

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ConvertFileToBytes convertit un fichier en série d'octets
func ConvertFileToBytes(file string) []byte {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

// ConvertBytesToFile convertit une série d'octets en fichier
func ConvertBytesToFile(name string, bytesArr []byte, perm int) {
	err := ioutil.WriteFile(name, bytesArr, os.FileMode(perm))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// GetByteSignature retourne la signature d'un packet
func GetByteSignature(packet []byte) [20]byte {
	return sha1.Sum(packet)
}

// GetFileSignature convertit un fichier en série d'octets et affiche la signature numérique du fichier
func GetFileSignature(file string) {
	data := ConvertFileToBytes(file)
	fmt.Printf("File signature : %x\n", sha1.Sum(data))
}

// GetFileByteSignature affiche la signature d'un fichier
func GetFileByteSignature(fileByte []byte) {
	log.Printf("File signature : %x\n\n", sha1.Sum(fileByte))
}
