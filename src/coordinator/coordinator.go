package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
)

var (
	peers *list.List
)

func main() {

	peers = list.New()

	utilities.Registration(peers, utilities.ServerPort, "coordinator")
	fmt.Println("Registration completed successfully!")

	for e := peers.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.ServerPort))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go handleConnection(connection)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	msg := new(tokenAsking.Message)
	dec := gob.NewDecoder(conn)
	err := dec.Decode(msg)
	utilities.CheckError(err, "error decoding message")

	if msg.MsgType == tokenAsking.Request {
		fmt.Println("ho ricevuto req")
	}

}
