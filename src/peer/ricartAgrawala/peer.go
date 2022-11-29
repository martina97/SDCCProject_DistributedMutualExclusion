package ricartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

type State string

const (
	Requesting State = "Requesting"
	CS         State = "CS"  //sto in sezione critica
	NCS        State = "NCS" //non in sezione critica
)

type RApeer struct {
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
	Num      lamport.ScalarClock
	lastReq  lamport.ScalarClock //timestamp del msg di richiesta
	state    State
	//Waiting  bool
	ChanStartTest chan bool

	PeerList *list.List //lista peer

	DeferSet *list.List
	replies  int //numero di risposte ricevute (inizializzato a 0)

}

func (p RApeer) GetMutex() sync.Mutex {
	return p.mutex
}

func NewRicartAgrawalaPeer(username string, ID int, address string, port string) *RApeer {
	peer := &RApeer{
		Username: username,
		ID:       ID,
		Address:  address,
		Port:     port,
		state:    NCS,
		DeferSet: list.New(),
		LogPath:  "/docker/node_volume/ricartAgrawala/peer_" + strconv.Itoa(ID) + ".log",
		//ChanRcvMsg = make(chan utilities.Message, utilities.MSG_BUFFERED_SIZE)
		//ChanSendMsg = make(chan *utilities.Message, utilities.MSG_BUFFERED_SIZE)
		ChanStartTest: make(chan bool, utilities.ChanSize),
	}
	peer.setInfos()
	return peer

}

func (m *RApeer) ToString() string {

	return fmt.Sprintf("myRapeer: {%s, num = %d, lastReq = %d, state = %s", m.Username, m.Num, m.lastReq, m.state+"}")
}

func (p *RApeer) setInfos() {
	utilities.CreateLog(p.LogPath, "[peer]") // in node_information.go

	f, err := os.OpenFile(p.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	_, err = f.WriteString("Initial timestamp of " + p.Username + " is " + strconv.Itoa(int(p.Num)))
	_, err = f.WriteString("\n")

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(f)

	lamport.StartSC(p.Num)

}
