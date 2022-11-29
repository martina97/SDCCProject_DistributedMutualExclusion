package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"fmt"
	"net"
)

var (
	peers *list.List
)

func main() {

	listener, err := net.Listen("tcp", ":1234")
	utilities.CheckError(err, "error listening")
	defer listener.Close()

	utilities.Registration(peers, utilities.ServerPort, "coordinator")
	fmt.Println("Registration completed successfully!")
	fmt.Println("peers == ", peers)
	for e := peers.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
