package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
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
	VC       VectorClock

	Coordinator Coordinator
}

func NewTokenAskingPeer(username string, ID int, address string, port string) *TokenPeer {
	peer := &TokenPeer{
		Username: username,
		ID:       ID,
		Address:  address,
		Port:     port,
		LogPath:  "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(ID) + ".log",
		HasToken: make(chan bool, utilities.CHAN_SIZE),
		VC:       make(map[string]int),
	}
	peer.setInfos()
	return peer
}

func (p *TokenPeer) setInfos() {
	StartVC(p.VC)
	utilities.CreateLog(p.LogPath, "[peer]")

	f, err := os.OpenFile(p.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	_, err = f.WriteString("Initial vector clock of " + p.Username + " is " + ToString(p.VC))
	_, err = f.WriteString("\n")

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(f)
}
