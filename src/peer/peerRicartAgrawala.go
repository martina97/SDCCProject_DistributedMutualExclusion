package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"log"
	"net"
	"sync"
)

type RA_peer struct {
	//info su nodo
	Username string //nome nodo
	ID       int    //id nodo
	Address  string //indirizzo nodo
	Port     string //porta nodo

	// utili per mutua esclusione
	mutex     sync.Mutex
	Timestamp utilities.TimeStamp
	fileLog   *log.Logger //file di ogni processo in cui scrivo info di quando accede alla sez critica
	Listener  net.Listener

	/*
		replies int //numero di risposte ricevute (inizializzato a 0)
		state string // Requesting, CS, NCS (inizializzato a NCS)
		queue *list.List // coda di richeste pendenti (inizialmente vuota)
		lastReq utilities.Timestamp //timestamp del msg di richiesta
		num utilities.Timestamp //clock logico scalare

	*/

	/*
		// algorithim
			shouldDefer     bool //è lo stato!!!!!!
			requestTS       msgp2.TimeStamp
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

func NewRicartAgrawalaPeer(username string, ID int, address string, port string) *RA_peer {
	return &RA_peer{Username: username, ID: ID, Address: address, Port: port}
}
