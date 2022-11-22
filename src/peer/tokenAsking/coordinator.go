package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
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

	PeerList *list.List //lista peer
	VC       utilities.VectorClock
	ReqList  *list.List
	HasToken bool
}

func NewCoordinator(username string, ID int, address string, port string, isCoord bool) *Coordinator {
	coordinator := &Coordinator{
		Username: username,
		ID:       ID,
		Address:  address,
		Port:     port,
		ReqList:  list.New(),
		LogPath:  "/docker/node_volume/tokenAsking/coordinator.log",
		//ChanRcvMsg = make(chan utilities.Message, utilities.MSG_BUFFERED_SIZE)
		//ChanSendMsg = make(chan *utilities.Message, utilities.MSG_BUFFERED_SIZE)
		VC:       make(map[string]int),
		HasToken: true,
	}

	utilities.StartVC2(coordinator.VC)
	if isCoord {
		coordinator.setInfos()
	}
	return coordinator

}

func (c *Coordinator) setInfos() {
	fmt.Println("sono in setInfos, logPAth == " + c.LogPath)
	/*
		c.VC = make(map[string]int)
		utilities.StartVC2(c.VC)

	*/
	c.ReqList.Init()
	utilities.CreateLog2(c.LogPath, "[coordinator]")

	f, err := os.OpenFile(c.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	date := time.Now().Format(utilities.DATE_FORMAT)
	_, err = f.WriteString("[" + date + "] : initial vector clock of coordinator is " + utilities.ToString(c.VC) + ".")
	_, err = f.WriteString("\n")
	date = time.Now().Format(utilities.DATE_FORMAT)
	_, err = f.WriteString("[" + date + "] : coordinator owns the token in starting up. ")
	_, err = f.WriteString("\n")

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(f)

}
