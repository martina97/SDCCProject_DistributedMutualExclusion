package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"log"
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
	mutex sync.Mutex

	//Waiting  bool
	HasToken bool

	PeerList *list.List //lista peer
	VC       VectorClock

	ChanStartTest chan bool
}

func NewTokenAskingPeer(username string, ID int, address string, port string) *TokenPeer {
	peer := &TokenPeer{
		Username:      username,
		ID:            ID,
		Address:       address,
		Port:          port,
		LogPath:       "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(ID) + ".log",
		HasToken:      false,
		VC:            make(map[string]int),
		ChanStartTest: make(chan bool, utilities.ChanSize),
	}
	peer.setInfos()
	return peer
}

func (p *TokenPeer) setInfos() {
	StartVC(p.VC)

	if verbose {
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

}
