package serverfunc

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

// SendPaquetWithFiability define fiability
func SendPaquetWithFiability(fiability float32) bool {
	if rand.Float32() <= fiability {
		return true
	}
	return false
}

// IsTimeOutError define if we have a timeout error
func IsTimeOutError(err error) int {
	if e, ok := err.(net.Error); ok && e.Timeout() {
		fmt.Fprintf(os.Stderr, "Timeout Error : %s\n", err.Error())
		return 1
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error : %s\n", err.Error())
		return -1
	}
	return 0
}

// NewClientConnexion etablish a connexion with a client
func NewClientConnexion(listener net.PacketConn) net.Addr {
	buff := make([]byte, 1000)

	listener.SetReadDeadline(time.Now().Add(5 * time.Minute))
	for {
		n, co, err := listener.ReadFrom(buff)
		if e, ok := err.(net.Error); ok && e.Timeout() {
			fmt.Fprintf(os.Stderr, "Timeout Error : %s\n", err.Error())
			return nil
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal Error : %s\n", err.Error())
			continue
		}
		if string(buff[:n]) == "Client - CONNEXION OK" {
			listener.WriteTo([]byte("Serveur - CONNEXION OK"), co)
			log.Println("Connexion established")
			return co
		}
	}
}
