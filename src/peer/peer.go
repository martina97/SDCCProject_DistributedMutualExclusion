package main

//è il main di ogni peer

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"bufio"
	"container/list"
	"fmt"
	"net"
	"os"
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
	myCoordinator tokenAsking.Coordinator
	algorithm     string

	//utili per test
	numSenders int
)

func main() {

	// Per fare la registrazione, il peer deve specificare il proprio nome (in questo modo nel file metto il nome del peer)
	peers = list.New()

	if utilities.Test {
		numSenders, algorithm = openTestMenu()
	}

	fmt.Println("Choose a username")
	in := bufio.NewReader(os.Stdin)
	// 2 peer non possono avere stesso username
	myUsername, _ = in.ReadString('\n')
	myUsername = strings.TrimSpace(myUsername)

	listener, err := net.Listen("tcp", ":1234")
	utilities.CheckError(err, "error listening")
	defer listener.Close()

	/* passo il result file a registration in modo che in esso vengono inserite
	le info del file!
	*/
	utilities.Registration(peers, utilities.ClientPort, myUsername)

	fmt.Println("Registration completed successfully!")
	//a questo punto tutti sanno quali sono gli altri peer

	setID()
	//open listen channel for messages
	//service on port 2345

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
		if myUsername == utilities.COORDINATOR {
			myCoordinator = *tokenAsking.NewCoordinator(myUsername, myID, myNode.Address, myNode.Port, true)
			myCoordinator.PeerList = peers
		} else {
			myTokenPeer = *tokenAsking.NewTokenAskingPeer(myUsername, myID, myNode.Address, myNode.Port)
			myTokenPeer.PeerList = peers
			for e := peers.Front(); e != nil; e = e.Next() {
				peer := e.Value.(utilities.NodeInfo)
				if peer.Username == utilities.COORDINATOR {
					myTokenPeer.Coordinator = *tokenAsking.NewCoordinator(peer.Username, peer.ID, peer.Address, peer.Port, false)
				}
			}
		}
	}
}
