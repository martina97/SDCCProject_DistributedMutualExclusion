package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
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

	//file di log
	LogPath string

	// utili per mutua esclusione
	mutex sync.Mutex

	VC            tokenAsking.VectorClock
	ReqList       *list.List
	HasToken      bool
	numTokenMsgs  int
	ChanStartTest chan bool
}

func NewCoordinator() *Coordinator {
	coordinator := &Coordinator{
		ReqList:       list.New(),
		LogPath:       "/docker/coordinator_volume/coordinator.log",
		VC:            make(map[string]int),
		HasToken:      true,
		numTokenMsgs:  0,
		ChanStartTest: make(chan bool, utilities.ChanSize),
	}

	tokenAsking.StartVC(coordinator.VC)
	coordinator.setInfos()
	return coordinator

}

func (c *Coordinator) setInfos() {
	c.ReqList.Init()
	utilities.CreateLog(c.LogPath, "[coordinator]")

	f, err := os.OpenFile(c.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	date := time.Now().Format(utilities.DateFormat)
	_, err = f.WriteString("[" + date + "] : initial vector clock of coordinator is " + tokenAsking.ToString(c.VC) + ".")
	_, err = f.WriteString("\n")
	date = time.Now().Format(utilities.DateFormat)
	_, err = f.WriteString("[" + date + "] : coordinator owns the token in starting up. ")
	_, err = f.WriteString("\n")

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(f)

}
