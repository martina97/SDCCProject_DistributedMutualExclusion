package main

//è il main di ogni peer

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	listNodes     []utilities.NodeInfo
	peers         *list.List
	myID          int
	myUsername    string
	allID         []int
	myNode        utilities.NodeInfo
	myLamportPeer lamport.LamportPeer
	myRApeer      ricartAgrawala.RApeer
	myTokenPeer   tokenAsking.TokenPeer
	algorithm     string

	//utili per test
	numSenders int
	verbose    bool
)

func main() {
	flag.BoolVar(&verbose, "v", utilities.Verbose, "use this flag to get verbose info on messages")
	flag.Parse()

	if verbose {
		fmt.Println("VERBOSE FLAG ON")
	}

	// Per fare la registrazione, il peer deve specificare il proprio nome (in questo modo nel file metto il nome del peer)
	peers = list.New()

	if utilities.Test {
		numSenders, algorithm = openTestMenu()
	}

	fmt.Println("Choose a username")
	in := bufio.NewReader(os.Stdin)
	myUsername, _ = in.ReadString('\n')
	myUsername = strings.TrimSpace(myUsername)

	listener, err := net.Listen("tcp", ":1234")
	utilities.CheckError(err, "error listening")
	defer listener.Close()

	utilities.Registration(peers, utilities.ClientPort, myUsername)

	fmt.Println("Registration completed successfully!")
	//a questo punto tutti sanno quali sono gli altri peer

	setID()

	go message_handler()

	if utilities.Test {
		//lancio i test
		runTest()
	} else {
		openMenu() //qui devo scegliere tra Lamport, Ricart-Agrawala e Token-Asking
	}
}

//Setto le variabili globali myNode, myID e allID
func setID() {
	//in peers ci sono tutti i peer
	for i := peers.Front(); i != nil; i = i.Next() {
		elem := i.Value.(utilities.NodeInfo)

		if elem.Username == myUsername {
			myNode = elem
			myID = elem.ID
			allID = append(allID, myID)
		} else {
			allID = append(allID, myID)
		}

	}
}

// Setto i peer in base all'algoritmo scelto
func setAlgorithmPeer() {
	switch algorithm {

	case "lamport":
		myLamportPeer = *lamport.NewLamportPeer(myUsername, myID, myNode.Address, myNode.Port)
		myLamportPeer.PeerList = peers

	case "ricartAgrawala":
		myRApeer = *ricartAgrawala.NewRicartAgrawalaPeer(myUsername, myID, myNode.Address, myNode.Port)
		myRApeer.PeerList = peers

	case "tokenAsking":

		myTokenPeer = *tokenAsking.NewTokenAskingPeer(myUsername, myID, myNode.Address, myNode.Port)
		myTokenPeer.PeerList = peers

	}
}

func message_handler() {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.ClientPort))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}

		switch algorithm {
		case "ricartAgrawala":
			go ricartAgrawala.HandleConnection(connection, &myRApeer)
		case "tokenAsking":

			go tokenAsking.HandleConnectionPeer(connection, &myTokenPeer)

		case "lamport":
			go lamport.HandleConnection(connection, &myLamportPeer)

		}

	}
}
