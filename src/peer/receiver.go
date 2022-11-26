package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	msgScaFile *os.File
)

func message_handler() {

	addr, err := net.ResolveTCPAddr("tcp", ":2345")
	if err != nil {
		fmt.Printf("Unable to resolve IP")
	}

	//listener, err := net.ListenTCP("tcp", ":"+strconv.Itoa(utilities.Client_port))
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	//defer listener.Close()

	/*
		//open file for save msg
		open_files()
		defer close_files()

	*/

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal("Accept fail")
		}

		// Enable Keepalives
		err = connection.SetKeepAlive(true)
		if err != nil {
			fmt.Printf("Unable to set keepalive - %s", err)
		}

		switch algorithm {
		case "ricartAgrawala":
			go ricartAgrawala.HandleConnection(connection, &myRApeer)
		case "tokenAsking":
			if myNode.Username == utilities.COORDINATOR {
				go tokenAsking.HandleConnectionCoordinator(connection, &myCoordinator)
			} else {
				go tokenAsking.HandleConnectionPeer(connection, &myTokenPeer)
			}
		case "lamport":
			go lamport.HandleConnection(connection, &myLamportPeer)

		}

		//go handleConnection(connection)
		//go handleConnectionCentralized(connection)
	}
}

/*
func open_files() {
	var err error
	msgScaFile, err = os.OpenFile(utilities.Peer_msg_sca_file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}
}

func close_files() {

	err := msgScaFile.Close()
	if err != nil {
		log.Fatal(err)
	}

}

*/
