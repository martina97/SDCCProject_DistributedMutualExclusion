package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	msgScaFile *os.File
)

func message_handler() {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.Client_port))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	//defer listener.Close()

	/*
		//open file for save msg
		open_files()
		defer close_files()

	*/

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
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
