package main

//è il main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"bufio"
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	listNodes  []utilities.NodeInfo
	peers      *list.List
	myID       int
	myUsername string
	allID      []int
	myNode     utilities.NodeInfo
	myRApeer   RApeer
	algorithm  string
	//lock   utilities.InfoLock
	//devo avere 3 peer (nodi), e su ogni peer viene eseguito un processo

	/*
		scalarMap utilities.MessageMap
		timeStamp utilities.Num //todo: mettere tutte queste var in una struttura per ogni processo

	*/
)

func main() {
	//listaNodi := []string{}
	var res utilities.Result_file //contiene stringhe del file di log client.txt

	//msg := []string{"messaggio1"}
	/*
		Per fare la registrazione, il peer deve specificare il proprio nome (in questo modo nel file metto il nome del peer)
	*/
	peers = list.New()

	fmt.Println("Choose a username")
	in := bufio.NewReader(os.Stdin)
	// 2 peer non possono avere stesso username
	myUsername, _ = in.ReadString('\n')
	myUsername = strings.TrimSpace(myUsername)

	listener, err := net.Listen("tcp", ":1234")
	utilities.Check_error(err)
	defer listener.Close()

	/* passo il result file a registration in modo che in esso vengono inserite
	le info del file!
	*/
	res = utilities.Registration(peers, utilities.Client_port, myUsername, listNodes)

	fmt.Println("PROVA DOPO REG")
	fmt.Println("res.PeerNum ====", res.PeerNum)
	fmt.Println("res.Peers == ", res.Peers, "\n\n")
	fmt.Println("NUMERO LIST PEER == ", peers.Len())
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
	setPeerUtils()
	//processo relativo al nodo che sto  considerando. il processo avrà info su
	// id nodo e indirizzo e porta nodo
	//fmt.Println("MY PEER =====", myNode.Username)

	//fmt.Println("sono il peer ", myUsername, "il mio id ===", myID)

	//startClocks()
	utilities.StartTS(myNode.TimeStamp)
	fmt.Println("START CLOCKS TERMINATO")
	fmt.Println("OPEN MENU")

	//open listen channel for messages
	//service on port 2345
	go message_handler()
	//go message_handler_centr()

	openMenu() //qui devo scegliere tra Lamport e Ricart Agrawala

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

		if elem.Username == myUsername {
			myNode = elem
			myID = elem.ID
			allID = append(allID, myID)
			//fmt.Println(" SONO ", myUsername, "IL MIO ID == ", myID)
		} else {
			allID = append(allID, myID)
		}

	}
}

func setPeerUtils() {
	// creo file "peer_ID.log"

	utilities.CreateLog(&myNode, "process_", strconv.Itoa(myNode.ID), "[peer]") // in nodeIdentification.go

	f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myNode.ID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	_, err = f.WriteString("Initial timestamp of process(" + strconv.Itoa(myID) + ") is " + strconv.Itoa(int(myNode.TimeStamp)))
	_, err = f.WriteString("\n")

	defer f.Close()

	/* todo: scommentare
	myNode.FileLog.SetOutput(f)
	myNode.FileLog.Println("infoProcess(" + strconv.Itoa(myNode.ID) + ") created.\n")

	*/

	//fmt.Println("logger ???? ", logger)

	//setto info sul processo in esecuzione sul peer
	//todo: serve?
	//NewProcess(&myNode)

	myNode.ChanRcvMsg = make(chan utilities.Message, utilities.MSG_BUFFERED_SIZE)
	myNode.ChanSendMsg = make(chan *utilities.Message, utilities.MSG_BUFFERED_SIZE)
	myNode.ChanAcquireLock = make(chan bool, utilities.CHAN_SIZE)
	myNode.SetDeferProSet(list.New())
	myNode.SetReplyProSet(list.New())

}

/*
func CreateLog(typeInfo string, id string, header string) *log.Logger {
	serverLogFile, err := os.Create("/docker/node_volume/" + typeInfo + id + ".log")
	if err != nil {
		log.Printf("unable to read file: %v", err)
	}
	serverLogFile.Close()
	/*
		newpath := filepath.Join(".", "log")
		os.MkdirAll(newpath, os.ModePerm)
		serverLogFile, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		// return log.New(serverLogFile, header, log.Lmicroseconds|log.Lshortfile)


	return log.New(serverLogFile, header, log.Lshortfile)
}

*/

func setAlgorithmPeer(algo string) {
	fmt.Println(" -------  sto in setAlgorithmPeer  -------")
	switch algo {
	case "ricart":
		myRApeer = *NewRicartAgrawalaPeer(myUsername, myID, myNode.Address, myNode.Port)
		fmt.Println("myRApeer ====", myRApeer)
		fmt.Println("myNode ====", myNode)
		utilities.StartTS(myRApeer.Num)
		fmt.Println("myRApeer.Num ==== ", myRApeer.Num)

	}

}
