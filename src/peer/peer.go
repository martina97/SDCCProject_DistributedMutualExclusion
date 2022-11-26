package main

//è il main

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
	utilities.Registration(peers, utilities.Client_port, myUsername, listNodes)

	//a questo punto tutti sanno quali sono gli altri peer

	setID()
	//open listen channel for messages
	//service on port 2345

	go message_handler()

	if utilities.Test {
		//lancio i test
		startTests()
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

func setAlgorithmPeer() {
	switch algorithm {
	case "lamport":
		myLamportPeer = *lamport.NewLamportPeer(myUsername, myID, myNode.Address, myNode.Port)
		//utilities.StartTS(myLamportPeer.Timestamp)
		myLamportPeer.PeerList = peers

	case "ricartAgrawala":
		myRApeer = *ricartAgrawala.NewRicartAgrawalaPeer(myUsername, myID, myNode.Address, myNode.Port)
		fmt.Println("myRApeer ====", myRApeer)
		fmt.Println("myNode ====", myNode)
		utilities.StartTS(myRApeer.Num)
		fmt.Println("myRApeer.Num ==== ", myRApeer.Num)

		//myRApeer.LogPath = "/docker/node_volume/ricartAgrawala/peer_" + strconv.Itoa(myRApeer.ID+1) + ".log"
		myRApeer.PeerList = peers
		fmt.Println("myRApeer.PeerList = ", myRApeer.PeerList)

	case "tokenAsking":
		if myUsername == utilities.COORDINATOR {
			myCoordinator = *tokenAsking.NewCoordinator(myUsername, myID, myNode.Address, myNode.Port, true)
			fmt.Println("myCoordinator ====", myCoordinator)
			fmt.Println("myNode ====", myNode)
			/*
				myCoordinator.VC = make(map[string]int)
				utilities.StartVC(myCoordinator.VC)

			*/
			myCoordinator.PeerList = peers
			for e := myCoordinator.PeerList.Front(); e != nil; e = e.Next() {
				peer := e.Value.(utilities.NodeInfo)
				if peer.Username != utilities.COORDINATOR {
					peer.LogPath = "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(peer.ID) + ".log"
				}
			}
			fmt.Println("myCoordinator.PeerList = ", myCoordinator.PeerList)
		} else {
			myTokenPeer = *tokenAsking.NewTokenAskingPeer(myUsername, myID, myNode.Address, myNode.Port)
			fmt.Println("myTokenPeer ====", myTokenPeer)
			fmt.Println("myNode ====", myNode)
			//myTokenPeer.VC = make(map[string]int)
			//utilities.StartVC(myTokenPeer.VC)
			fmt.Println("myTokenPeer.VC =", myTokenPeer.VC)
			myTokenPeer.PeerList = peers
			fmt.Println("myTokenPeer.PeerList = ", myTokenPeer.PeerList)
			//fmt.Println("toString 2 ----", (myTokenPeer.VC).ToString2())
			for e := peers.Front(); e != nil; e = e.Next() {
				fmt.Println("e ==", e)
				fmt.Println("e.Value ==", e.Value)
				peer := e.Value.(utilities.NodeInfo)
				if peer.Username == utilities.COORDINATOR {
					fmt.Println("il coordinatore è = ", peer.Username)
					myTokenPeer.Coordinator = *tokenAsking.NewCoordinator(peer.Username, peer.ID, peer.Address, peer.Port, false)

					/*
						utilities.StartVC(myTokenPeer.Coordinator.VC)

					*/
				}
			}
		}

	}

}
