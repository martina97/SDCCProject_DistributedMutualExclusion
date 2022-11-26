package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"log"
	"os"
	"sync"
	"time"
)

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
)

type Coordinator struct {
	//info su nodo
	Username string //nome nodo
	ID       int    //id nodo
	Address  string //indirizzo nodo
	Port     string //porta nodo

	//file di log
	LogPath string

	// utili per mutua esclusione
	mutex sync.Mutex

	PeerList     *list.List //lista peer
	VC           utilities.VectorClock
	ReqList      *list.List
	HasToken     bool
	numTokenMsgs int
}

func NewCoordinator(username string, ID int, address string, port string, isCoord bool) *Coordinator {
	coordinator := &Coordinator{
		Username:     username,
		ID:           ID,
		Address:      address,
		Port:         port,
		ReqList:      list.New(),
		LogPath:      "/docker/node_volume/tokenAsking/coordinator.log",
		VC:           make(map[string]int),
		HasToken:     true,
		numTokenMsgs: 0,
	}

	utilities.StartVC(coordinator.VC)
	if isCoord {
		coordinator.setInfos()
	}
	return coordinator

}

func (c *Coordinator) setInfos() {
	var err error

	c.ReqList.Init()
	utilities.CreateLog2(c.LogPath, "[coordinator]")

	f := openFile(true)
	date := time.Now().Format(utilities.DATE_FORMAT)
	_, err = f.WriteString("[" + date + "] : initial vector clock of coordinator is " + utilities.ToString(c.VC) + ".")
	_, err = f.WriteString("\n")
	date = time.Now().Format(utilities.DATE_FORMAT)
	_, err = f.WriteString("[" + date + "] : coordinator owns the token in starting up. ")
	_, err = f.WriteString("\n")

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(f)

}
