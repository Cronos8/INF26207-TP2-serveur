package filebyte

import (
	"fmt"
	"io/ioutil"
	"crypto/sha1"
	"os"
	"log"
)

// ConvertFileToBytes convert file to byte
func ConvertFileToBytes(file string) []byte {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

// ConvertBytesToFile convert byte to file
func ConvertBytesToFile(name string, bytesArr []byte, perm int) {
	err := ioutil.WriteFile(name, bytesArr, os.FileMode(perm))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// GetByteSignature
func GetByteSignature(packet []byte) [20]byte {
	return sha1.Sum(packet)
}

// GetFileSignature
func GetFileSignature(file string) {
	data := ConvertFileToBytes(file)
	fmt.Printf("File signature : %x\n", sha1.Sum(data))
}

// GetFileByteSignature
func GetFileByteSignature(fileByte []byte) {
	log.Printf("File signature : %x\n", sha1.Sum(fileByte))
}
