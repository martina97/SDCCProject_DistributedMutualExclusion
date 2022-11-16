package RicartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
	"log"
	"net"
	"strconv"
)

func Message_handler() {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.Client_port))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go HandleConnection(connection)
		//go handleConnectionCentralized(connection)
	}
}

//Save message
func HandleConnection(conn net.Conn) error {

	fmt.Println("sto in handleConnection dentro RicartAgrawala package")

	return nil
}
