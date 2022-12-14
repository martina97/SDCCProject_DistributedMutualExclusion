package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"

	//"provaMutualExLamport/src/utilities"
	"strconv"
)

type Register struct{}

//var info utilities.nodeInfo //info e' struttura nodo
//var listNodes [3]utilities.nodeInfo

func main() {
	var connect_num int
	utility := new(utilities.Utility)

	server := rpc.NewServer()
	//register method
	err := server.RegisterName("Utility", utility)
	if err != nil {
		log.Fatal("Format of service Utility is not correct: ", err)
	}

	// Create file to maintain ip address and number port of all registered nodes
	f, err := os.Create("/docker/register_volume/nodes.txt")
	if err != nil {
		log.Printf("unable to read file: %v", err)
	}
	f.Close()

	port := 4321
	// listen for a connection
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Error in listening:", err)
	}
	// Close the listener whenever we stop
	defer listener.Close()

	log.Printf("RPC server on port %d", port)

	go server.Accept(listener)

	//Wait connection
	for connect_num < utilities.MAXPEERS { //todo: mettere 3 , anche sotto
		ch := <-utilities.Connection
		if ch == true {
			connect_num++
		}
	}

	log.Printf("Max Number of Connection reached up")

	utilities.Wg.Add(-utilities.MAXPEERS) //quando tutti e 3 sono registrati si ritorna a #src

	//send client a responce for max number of peer registered
	time.Sleep(time.Minute)

	/*
		for {
			//TODO after registration this peer must be off ??
		}

	*/

}
