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
	listNodes []utilities.NodeInfo
	peers     *list.List
	//myID          int
	myUsername string
	//allID         []int
	myNode        utilities.NodeInfo
	myLamportPeer lamport.LamportPeer
	myRApeer      ricartAgrawala.RApeer
	myTokenPeer   tokenAsking.TokenPeer
	myCoordinator tokenAsking.Coordinator
	algorithm     string

	//utile per test
	numSenders int
)

func main() {
	//listaNodi := []string{}
	var res utilities.Result_file //contiene stringhe del file di log client.txt

	//msg := []string{"messaggio1"}
	/*
		Per fare la registrazione, il peer deve specificare il proprio nome (in questo modo nel file metto il nome del peer)
	*/
	peers = list.New()

	if utilities.Test {
		numSenders, algorithm = openTestMenu()
	}

	fmt.Println("Choose a username")
	in := bufio.NewReader(os.Stdin)
	// 2 peer non possono avere stesso username
	myNode.Username, _ = in.ReadString('\n')
	myNode.Username = strings.TrimSpace(myNode.Username)

	listener, err := net.Listen("tcp", ":1234")
	utilities.CheckError(err)
	defer listener.Close()

	/* passo il result file a registration in modo che in esso vengono inserite
	le info del file!
	*/
	res = utilities.Registration(peers, utilities.Client_port, myNode.Username, listNodes)

	fmt.Println("PROVA DOPO REG")
	fmt.Println("res.PeerNum ====", res.PeerNum)
	fmt.Println("res.Peers == ", res.Peers, "\n\n")
	fmt.Println("NUMERO LIST PEER == ", peers.Len())
	fmt.Println("LIST PEER ===== ", peers)
	fmt.Println("*LIST PEER ===== ", *peers)
	fmt.Println("&LIST PEER ===== ", &peers)

	for e := peers.Front(); e != nil; e = e.Next() {
		fmt.Println("e ==", e)
	}

	/*
		for e := peers.Front(); e != nil; e = e.Next() {
			// do something with e.Value
			fmt.Println("PROVA STAMPA PEER1", reflect.TypeOf(e.Value))

			fmt.Println("PROVA STAMPA PEER1", e.Value)

		}

	*/

	//a questo punto tutti sanno quali sono gli altri peer
	// prova a mandare dal peer marti un messaggio agli altri 2
	// TODO: per ogni peer faccio lock ??!???
	//devo prendermi ID del nodo e la porta del nodo!
	setID()
	//processo relativo al nodo che sto  considerando. il processo avrà info su
	// id nodo e indirizzo e porta nodo
	//fmt.Println("MY PEER =====", myNode.Username)

	//fmt.Println("sono il peer ", myUsername, "il mio id ===", myID)

	//startClocks()

	fmt.Println("START CLOCKS TERMINATO")
	fmt.Println("OPEN MENU")

	//open listen channel for messages
	//service on port 2345

	go message_handler()

	//go message_handler()
	//go message_handler_centr()

	if utilities.Test {
		//lancio i test
		//p0 sceglie il test (x semplicità 1 solo lo sceglie)
		//if myID == 0 {
		//}
		startTests()

	} else {
		openMenu() //qui devo scegliere tra Lamport e Ricart Agrawala

	}

	// creo file "peer_ID.log"

	//creo nuovo processo in esecuzione sul peer
	//p, err := NewProcess(myNode)

	//TODO: scommentare
}

func setID() {
	/* devo settare la variabile globale myID per sapere qual è ID del peer
	che dovra' fare determinate azioni, e devo creare una lista che ha tutti gli
	ID degli altri peer
	*/
	//ora setto la variabile globale myNode
	//scorro peers, che e' *list, in cui ci sono i peer
	for i := peers.Front(); i != nil; i = i.Next() {
		//fmt.Println("PROVA SET ID ", i.Value.(utilities.nodeInfo)) //i.value e' interface{}
		elem := i.Value.(utilities.NodeInfo)
		//fmt.Println(" myUsername ==== ", myUsername)
		//fmt.Println(" elem.Username ==== ", elem.Username)

		if elem.Username == myNode.Username {
			myNode = elem
			myNode.ID = elem.ID
			//allID = append(allID, myID)
			//fmt.Println(" SONO ", myUsername, "IL MIO ID == ", myID)
		} else {
			//allID = append(allID, myID)
		}

	}
}

func setAlgorithmPeer() {
	fmt.Println(" -------  sto in setAlgorithmPeer  -------")
	switch algorithm {
	case "lamport":
		myLamportPeer = *lamport.NewLamportPeer(myNode.Username, myNode.ID, myNode.Address, myNode.Port)
		fmt.Println("myLamportPeer ====", myLamportPeer)
		fmt.Println("myNode ====", myNode)
		utilities.StartTS(myLamportPeer.Timestamp)
		myLamportPeer.PeerList = peers
		fmt.Println("myLamportPeer.PeerList = ", myLamportPeer.PeerList)

	case "ricartAgrawala":
		myRApeer = *ricartAgrawala.NewRicartAgrawalaPeer(myNode.Username, myNode.ID, myNode.Address, myNode.Port)
		fmt.Println("myRApeer ====", myRApeer)
		fmt.Println("myNode ====", myNode)
		utilities.StartTS(myRApeer.Num)
		fmt.Println("myRApeer.Num ==== ", myRApeer.Num)

		//myRApeer.LogPath = "/docker/node_volume/ricartAgrawala/peer_" + strconv.Itoa(myRApeer.ID+1) + ".log"
		myRApeer.PeerList = peers
		fmt.Println("myRApeer.PeerList = ", myRApeer.PeerList)

	case "tokenAsking":
		if myNode.Username == utilities.COORDINATOR {
			myCoordinator = *tokenAsking.NewCoordinator(myNode.Username, myNode.ID, myNode.Address, myNode.Port, true)
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
			myTokenPeer = *tokenAsking.NewTokenAskingPeer(myNode.Username, myNode.ID, myNode.Address, myNode.Port)
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

				}
			}
		}

	}

}
