package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

type TokenPeer struct {
	//info su nodo
	Username string //nome nodo
	ID       int    //id nodo
	Address  string //indirizzo nodo
	Port     string //porta nodo

	//file di log
	LogPath string

	// utili per mutua esclusione
	mutex    sync.Mutex
	fileLog  *log.Logger //file di ogni processo in cui scrivo info di quando accede alla sez critica
	Listener net.Listener

	//Waiting  bool
	HasToken chan bool

	PeerList *list.List //lista peer
	VC       utilities.VectorClock

	Coordinator Coordinator
}

func NewTokenAskingPeer(username string, ID int, address string, port string) *TokenPeer {
	peer := &TokenPeer{
		Username: username,
		ID:       ID,
		Address:  address,
		Port:     port,
		LogPath:  "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(ID) + ".log",
		//ChanRcvMsg = make(chan utilities.Message, utilities.MSG_BUFFERED_SIZE)
		//ChanSendMsg = make(chan *utilities.Message, utilities.MSG_BUFFERED_SIZE)
		HasToken: make(chan bool, utilities.CHAN_SIZE),
		VC:       make(map[string]int),
	}
	peer.setInfos()
	return peer

}

func (p *TokenPeer) setInfos() {
	fmt.Println("sono in setInfos, logPAth == " + p.LogPath)

	utilities.StartVC2(p.VC)

	utilities.CreateLog2(p.LogPath, "[peer]") // in nodeIdentification.go

	f, err := os.OpenFile(p.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//_, err = f.WriteString("Initial timestamp of " + p.Username + " is " + strconv.Itoa(int(p.Num)))
	_, err = f.WriteString("Initial vector clock of " + p.Username + " is " + utilities.ToString(p.VC))

	_, err = f.WriteString("\n")

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(f)

	/* todo: scommentare
	myNode.FileLog.SetOutput(f)
	myNode.FileLog.Println("infoProcess(" + strconv.Itoa(myNode.ID) + ") created.\n")

	*/

	//fmt.Println("logger ???? ", logger)

	//setto info sul processo in esecuzione sul peer
	//todo: serve?
	//NewProcess(&myNode)

}
