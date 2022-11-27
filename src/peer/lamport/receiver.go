package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
)

func HandleConnection(conn net.Conn, peer *LamportPeer) {
	if myPeer.Username == "" {
		myPeer = *peer
	}

	// read msg and save on file
	defer conn.Close()
	msg := new(Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)

	//ogni volta che ricevo un msg devo aggiornare TS
	//aggiorno timestamp
	//ogni peer ha il suo clock scalare, e' var globale come myNode e myID

	//time.Sleep(time.Minute / 2) //PRIMA DI AUMENTARE TS METTO SLEEP COSI PROVO A INVIARE 2 REQ INSIEME E VEDO CHE SUCCEDE

	//mutex := lock.GetMutex()

	if msg.MsgType == Request {
		UpdateTS(&myPeer.Timestamp, &msg.TS)
		/*
			quando ricevo una richiesta da un processo devo decidere se mandare ACK al processo oppure se voglio entrare in CS
		*/
		myPeer.mutex.Lock()
		WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp)

		//metto msg in mappa
		AppendHashMap(myPeer.ScalarMap, *msg)

		date := time.Now().Format(utilities.DateFormat)
		replyMsg := NewReply(myPeer.Username, msg.Sender, date, myPeer.Timestamp)
		sendReply(replyMsg)
		myPeer.mutex.Unlock()
	}

	if msg.MsgType == Reply {
		myPeer.mutex.Lock()

		//utilities.WriteMsgToFile(&myNode, "Receive", *msg, 0, myNode.TimeStamp)
		WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp)

		//aggiungo a replyProSet il msg
		myPeer.replySet.PushBack(msg)
		//check ack
		//checkAcks() //controllo se ho ricevuto 2 msg reply, se si posso entrare in CS prendendo 1 elem nella lista
		// e controllando che id sia il mio, se e' il mio entro altrimenti no
		myPeer.mutex.Unlock()
	} else if msg.MsgType == Release {
		myPeer.mutex.Lock()

		WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp)

		RemoveFirstElementMap(myPeer.ScalarMap)
		//checkAcks()
		myPeer.numRelease++
		if myPeer.numRelease == numSender {
			myPeer.StartTest <- true
		}
		myPeer.mutex.Unlock()

	}

}

func checkAcks() {
	fmt.Println("sto in checkAcks")

	//todo: quando azzero lista ReplyProSet ?????
	//date := time.Now().Format("15:04:05.000")

	for !(myPeer.replySet.Len() == myPeer.PeerList.Len()-1 && len(myPeer.ScalarMap) > 0) {
		fmt.Println("sto in checkAcks dentro for")

		time.Sleep(time.Second * 5)
	}
	fmt.Println("sto in checkAcks fuori for")

	//if myPeer.replySet.Len() == myPeer.PeerList.Len()-1 && len(myPeer.ScalarMap) > 0 {

	//prendo il primo mess nella mappa per vedere se è il mio, ossia guardo ID sender
	msg := GetFirstElementMap(myPeer.ScalarMap)

	if msg.Sender == myPeer.Username {
		myPeer.mutex.Lock()
		//il primo msg in lista è il mio, quindi posso accedere in CS
		//myPeer.Waiting = false
		//myPeer.ChanAcquireLock <- true

		utilities.WriteInfosToFile("receives all peer reply messages successfully.", myPeer.LogPath, myPeer.Username)

		//ho ricevuto tutti msg reply, ora entro in cs
		date := time.Now().Format(utilities.DateFormat)

		utilities.WriteInfosToFile("enters the critical section at "+date+".", myPeer.LogPath, myPeer.Username)

		time.Sleep(time.Minute / 2)
		date = time.Now().Format(utilities.DateFormat)

		utilities.WriteInfosToFile("exits the critical section at "+date+".", myPeer.LogPath, myPeer.Username)

		//lascio CS e mando msg release a tutti
		fmt.Println("devo mandare release!!!!!")
		sendRelease()
		myPeer.ChanAcquireLock <- true
		fmt.Println("chan true!!!!")
		myPeer.mutex.Unlock()
	}

}
