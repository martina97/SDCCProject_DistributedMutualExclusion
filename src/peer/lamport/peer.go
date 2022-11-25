package lamport

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

type LamportPeer struct {
	//info su nodo
	Username string //nome nodo
	ID       int    //id nodo
	Address  string //indirizzo nodo
	Port     string //porta nodo

	//file di log
	LogPath string

	// utili per mutua esclusione
	mutex     sync.Mutex
	Timestamp utilities.TimeStamp
	fileLog   *log.Logger //file di ogni processo in cui scrivo info di quando accede alla sez critica
	Listener  net.Listener

	//Waiting  bool
	Waiting         bool //serve a vedere se chi ha mandato msg request e' in attesa di tutti i msg reply
	ChanRcvMsg      chan utilities.Message
	ChanSendMsg     chan *utilities.Message
	ChanAcquireLock chan bool
	StartTest       chan bool
	//replyProSet     *list.List // then Message.Sender is the key.
	deferProSet *list.List // then Message.Sender is the key.

	ScalarMap utilities.MessageMap

	PeerList *list.List //lista peer

	DeferSet *list.List
	replySet *list.List

	numRelease int
}

func NewLamportPeer(username string, ID int, address string, port string) *LamportPeer {
	peer := &LamportPeer{
		Username: username,
		ID:       ID,
		Address:  address,
		Port:     port,
		DeferSet: list.New(),
		replySet: list.New(),
		LogPath:  "/docker/node_volume/lamport/peer_" + strconv.Itoa(ID) + ".log",
		//ChanRcvMsg = make(chan utilities.Message, utilities.MSG_BUFFERED_SIZE)
		//ChanSendMsg = make(chan *utilities.Message, utilities.MSG_BUFFERED_SIZE)
		ChanAcquireLock: make(chan bool, utilities.CHAN_SIZE),
		ChanRcvMsg:      make(chan utilities.Message, utilities.MSG_BUFFERED_SIZE),
		ChanSendMsg:     make(chan *utilities.Message, utilities.MSG_BUFFERED_SIZE),
		StartTest:       make(chan bool, utilities.CHAN_SIZE),
		deferProSet:     list.New(),
		ScalarMap:       utilities.MessageMap{},
		numRelease:      0,
	}
	peer.setInfos()
	return peer

}

func (p *LamportPeer) setInfos() {
	fmt.Println("sono in setInfos, logPAth == " + p.LogPath)
	utilities.CreateLog2(p.LogPath, "[peer]") // in nodeIdentification.go

	f, err := os.OpenFile(p.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	_, err = f.WriteString("Initial timestamp of " + p.Username + " is " + strconv.Itoa(int(p.Timestamp)))
	_, err = f.WriteString("\n")

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(f)

}
