package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
	"log"
	"net"
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
	logPath string;

	// utili per mutua esclusione
	mutex    sync.Mutex
	fileLog  *log.Logger //file di ogni processo in cui scrivo info di quando accede alla sez critica
	Listener net.Listener
	Num      utilities.TimeStamp
	lastReq  utilities.TimeStamp //timestamp del msg di richiesta
	state    State

	/*
		replies int //numero di risposte ricevute (inizializzato a 0)
		state string // Requesting, CS, NCS (inizializzato a NCS)
		queue *list.List // coda di richeste pendenti (inizialmente vuota)

		num utilities.Timestamp //clock logico scalare

	*/

	/*
		// algorithim
			shouldDefer     bool //è lo stato!!!!!!
			requestTS       msgp2.Num
			replyProSet     *list.List
			deferProSet     *list.List
			chanRcvMsg      chan msgp2.Message
			chanSendMsg     chan *msgp2.Message
			chanAcquireLock chan bool
			logger          *log.Logger

			// process handler
			p *process
			// sata info
			readCnt  int
			writeCnt int
	*/
}

func (p RApeer) GetMutex() sync.Mutex {
	return p.mutex
}

func NewRicartAgrawalaPeer(username string, ID int, address string, port string) *RApeer {
	return &RApeer{Username: username, ID: ID, Address: address, Port: port, state: NCS}
}

func (m *RApeer) ToString() string {

	return fmt.Sprintf("myRapeer: {%s, num = %d, lastReq = %d, state = %s", m.Username, m.Num, m.lastReq, m.state+"}")
}
