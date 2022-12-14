package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type MessageMap map[ScalarClock][]Message

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
	Timestamp ScalarClock

	//Waiting  bool
	Waiting       bool //serve a vedere se chi ha mandato msg request e' in attesa di tutti i msg reply
	ChanStartTest chan bool
	StartTest     chan bool
	ScalarMap     MessageMap
	replySet      *list.List //lista in cui metto i msg di reply

	PeerList *list.List //lista peer

	numMsgsTest int
	numRelease  int
}

func NewLamportPeer(username string, ID int, address string, port string) *LamportPeer {
	peer := &LamportPeer{
		Username:      username,
		ID:            ID,
		Address:       address,
		Port:          port,
		replySet:      list.New(),
		LogPath:       "/docker/node_volume/lamport/peer_" + strconv.Itoa(ID) + ".log",
		ChanStartTest: make(chan bool, utilities.ChanSize),
		StartTest:     make(chan bool, utilities.ChanSize),
		ScalarMap:     MessageMap{},
		numRelease:    0,
		numMsgsTest:   0,
	}
	peer.setInfos()
	return peer
}

func (p *LamportPeer) setInfos() {

	if verbose {
		utilities.CreateLog(p.LogPath, "[peer]") // in node_information.go

		f, err := os.OpenFile(p.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		_, err = f.WriteString("Initial logical scalar clock of " + p.Username + " is " + strconv.Itoa(int(p.Timestamp)))
		_, err = f.WriteString("\n")

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatalf("error closing file: %v", err)
			}
		}(f)
	}

	StartSC(p.Timestamp)
}

// test inserts several ints into an MessageHeap, checks the minimum,
// and removes them in order of priority.

func AppendHashMap(map1 MessageMap, message Message) {
	var listMsg []Message

	_, ok := map1[message.TS] //controllo se nella mappa c'?? la chiave message.TS

	if ok == true {
		//la chiave ?? presente
		//prendo elementi value relativi a quella chiave e faccio controllo
		if len(map1[message.TS]) == 1 { //c'e un solo valore nella lista relativa al TS (1 solo msg con quel TS)
			if message.Sender < map1[message.TS][0].Sender {
				map1[message.TS] = append([]Message{message}, map1[message.TS]...) //inserisco il msg all'inizio dello slice
				/*
					NB: data := []string{"A", "B", "C", "D"}
					METTO ALLA FINE: data = append(data, "prova")	--> [A B C D prova]
					METTO ALL'INIZIO: data = append([]string{"prova"}, data...) --> [prova A B C D]
				*/
			}
		} else {
			for i := 1; i < len(map1[message.TS]); i++ {
				if map1[message.TS][i-1].Sender < message.Sender && message.Sender < map1[message.TS][i].Sender {
					//fmt.Println("IL MSG STA TRA I 2")
					// devo inserire il msg tra i-1 e i
					map1[message.TS] = append(map1[message.TS], message) // msg ora e' in posiz len(msg1)

					copy(map1[message.TS][i+1:], map1[message.TS][i:])
					map1[message.TS][i] = message
					/*
						se ho slice : arr =[1 3 5] e voglio aggiungere il 2 tra 1 e 3:
						1. metto il 2 alla fine --> arr = [1 3 5 2]
						2. arr[2:] == [5 2] e arr[1:] == [3 5 2], con copy copio [3 5 2] in [5 2], ottenendo arr = [1 3 3 5]
						3. poi dico che arr[1] = 2 --> arr = [1 2 3 5]
					*/

					break
				}
				if map1[message.TS][len(map1[message.TS])-1].Sender < message.Sender {
					map1[message.TS] = append(map1[message.TS], message) // metto msg alla fine
					break
				}
			}
		}

	} else { // nella mappa non c'e quella chiave
		map1[message.TS] = append(listMsg, message)
	}
}

func GetFirstElementMap(mapMsg MessageMap) Message {
	var message Message
	for _, element := range mapMsg {
		//fmt.Println("Key:", key, "=>", "Element:", element)
		message = element[0]
		break
	}
	//fmt.Println("GetFirstElementMap ------", message)
	return message

}

func RemoveFirstElementMap(mapMsg MessageMap) {

	for key, _ := range mapMsg {
		mapMsg[key] = mapMsg[key][1:]
		if len(mapMsg[key]) == 0 { //se non ci sono piu msg con quel TS, ossia la lista di msg per quel TS (key) e' vuota
			delete(mapMsg, key)

			break
		}
	}

	fmt.Println("mappa == ", mapMsg)

}
