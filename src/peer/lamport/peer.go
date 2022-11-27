package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"log"
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
	LogPath string //file di ogni processo in cui scrivo info di quando accede alla sez critica

	// utili per mutua esclusione
	mutex     sync.Mutex
	Timestamp utilities.TimeStamp

	//Waiting  bool
	Waiting         bool //serve a vedere se chi ha mandato msg request e' in attesa di tutti i msg reply
	ChanAcquireLock chan bool
	StartTest       chan bool
	ScalarMap       utilities.MessageMap
	replySet        *list.List //lista in cui metto i msg di reply

	PeerList *list.List //lista peer

	numRelease int
}

func NewLamportPeer(username string, ID int, address string, port string) *LamportPeer {
	peer := &LamportPeer{
		Username:        username,
		ID:              ID,
		Address:         address,
		Port:            port,
		replySet:        list.New(),
		LogPath:         "/docker/node_volume/lamport/peer_" + strconv.Itoa(ID) + ".log",
		ChanAcquireLock: make(chan bool, utilities.CHAN_SIZE),
		StartTest:       make(chan bool, utilities.CHAN_SIZE),
		ScalarMap:       utilities.MessageMap{},
		numRelease:      0,
	}
	peer.setInfos()
	return peer
}

func (p *LamportPeer) setInfos() {
	utilities.CreateLog(p.LogPath, "[peer]") // in nodeIdentification.go

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

	utilities.StartTS(p.Timestamp)
}
