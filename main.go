package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Cronos8/INF26207-TP2-serveur/serverfilebyte"
	"github.com/Cronos8/INF26207-TP2-serveur/serverfunc"
	"github.com/Cronos8/INF26207-TP2-serveur/serverpacket"
)

// Vérifie si le nombre d'arguments reçu est conforme
func checkArguments(args []string) int {
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr filepath\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s 127.0.0.1:22222 testfiles/alpagas.jpeg\n", os.Args[0])
		os.Exit(1)
		return -1
	}
	return 0
}

func main() {

	hpacket := serverpacket.HeaderPacket{HeaderIP: nil, HeaderPort: 0, HeaderNbPacket: 0}

	// Vérification du nombre de paramètre envoyé
	if checkArguments(os.Args) != 0 {
		return
	}

	// Vérifie si l'adresse IPv4 du serveur est conforme
	udpaddr, err := net.ResolveUDPAddr("udp4", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid IPv4 address...\n")
		os.Exit(1)
	}

	hpacket.HeaderIP = udpaddr.IP
	hpacket.HeaderPort = int32(udpaddr.Port)

	// Conversion du fichier en suite d'octets
	fByte := serverfilebyte.ConvertFileToBytes(os.Args[2])
	if fByte == nil {
		os.Exit(1)
	}

	// Création d'un listener à partir de l'adresse IPv4 reçu en paramètre
	listener, err := net.ListenPacket("udp4", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s\n", err.Error())
		os.Exit(1)
	}

	// Fermeture du listener à la fin
	defer listener.Close()

	fmt.Println()
	fmt.Println("------------------------------------------------")
	fmt.Printf("Listening on : %s\n", os.Args[1])
	fmt.Println("------------------------------------------------")
	fmt.Println()

	i := 0
	size := 1000
	buffRead := make([]byte, 1000)
	buffWrite := make([]byte, 1024)
	lastbuff := make([]byte, 1024)
	terminated := false

	var fiability float32 = 0.95
	var nbPacket uint64 = 0

	// Établissement de la connexion au client
	conn := serverfunc.NewClientConnexion(listener)
	if conn != nil {
		// Envoi le nom du fichier
		serverfunc.SendFileName(listener, conn, os.Args[2])
		// Mise en place d'un compteur de 3 secondes d'attente de recéption d'un méssage
		listener.SetReadDeadline(time.Now().Add(3 * time.Second))
		for terminated == false {
			// Lecture
			n, _, err := listener.ReadFrom(buffRead)
			switch e := serverfunc.IsTimeOutError(err); e {
			case 1: // Si le délai de 3 secondes est dépassé

				// Renvoie du dernier paquet
				_, err2 := listener.WriteTo(lastbuff, conn)
				if err2 != nil {
					fmt.Println(err2)
				}

				// Affichage du paquet
				serverpacket.PrintMessage("PACKET RE-SEND", serverpacket.CyanColor, conn.String())
				serverpacket.PrintPacket(lastbuff)
				fmt.Println()

				// Mise à jour du compteur
				listener.SetReadDeadline(time.Now().Add(3 * time.Second))
				break

			case -1: // Si une erreur est survenue à la lecture
				fmt.Println(e)
				break

			case 0: // Si le message à bien été reçu

				// Affichage du message reçu
				serverpacket.PrintMessage("MESSAGE RECEIVE", serverpacket.GreenColor, conn.String())
				fmt.Printf("Content : %s\n", string(buffRead[:n]))

				// Envoi du paquet au client avec une fiabilité de 95%
				if serverfunc.SendPaquetWithFiability(fiability) == true {

					// Si le serveur reçoit l'accusé de réception du client
					if (string(buffRead[:n]) == "PACKET RECEIVE") || (string(buffRead[:n]) == "READY TO RECEIVE") {

						// Si il s'agit du dernier paquet à transmettre
						if i+size > len(fByte) {

							// Encapsulation du dernier paquet à envoyer au client
							buffWrite = serverpacket.EncapPacket(hpacket, fByte[i:])

							// Envoi du paquet
							listener.WriteTo(buffWrite, conn)

							// Affichage du paquet envoyé
							serverpacket.PrintMessage("LAST PACKET", serverpacket.PurpleColor, conn.String())
							serverpacket.PrintPacket(buffWrite)
							fmt.Println()

							// Envoi du message de fin de transmission
							listener.WriteTo([]byte("END"), conn)

							// Affichage de le signature numérique du fichier transmis
							serverfilebyte.GetFileByteSignature(fByte)
							terminated = true
							break
						}

						// Encapsulation du paquet à envoyer au client
						buffWrite = serverpacket.EncapPacket(hpacket, fByte[i:i+size])

						// Envoi du paquet
						listener.WriteTo(buffWrite, conn)

						// Affichage du pquet envoyé
						serverpacket.PrintMessage("PACKET SEND", serverpacket.BlueColor, conn.String())
						serverpacket.PrintPacket(buffWrite)
						fmt.Println()

						nbPacket++
						hpacket.HeaderNbPacket = nbPacket
					}

					i = i + size

					// copie du dernier paquet envoyé
					lastbuff = buffWrite

					// Mise à jour du compteur
					listener.SetReadDeadline(time.Now().Add(3 * time.Second))

				} else {
					serverpacket.PrintMessage("FIABILITY ERROR", serverpacket.RedColor, conn.String())
				}
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Server timeout ...\n")
		os.Exit(1)
	}
}
