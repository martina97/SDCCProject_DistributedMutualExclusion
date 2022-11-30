package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var verbose bool

func HandleConnection(conn net.Conn, peer *LamportPeer) {
	flag.BoolVar(&verbose, "v", utilities.Verbose, "use this flag to get verbose info on messages")
	flag.Parse()

	if verbose {
		fmt.Println("VERBOSE FLAG ON")
	}

	if myPeer.Username == "" {
		myPeer = *peer
	}

	// read msg and save on file
	defer conn.Close()
	msg := new(Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)

	if msg.MsgType == Request {
		UpdateSC(&myPeer.Timestamp, &msg.TS)

		myPeer.mutex.Lock()
		if verbose {
			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp)
		}

		//metto msg in mappa
		AppendHashMap(myPeer.ScalarMap, *msg)

		date := time.Now().Format(utilities.DateFormat)
		replyMsg := NewReply(myPeer.Username, msg.Sender, date, myPeer.Timestamp)
		err := sendReply(replyMsg)
		if err != nil {
			log.Fatalf("error sending reply %v", err)
		}
		myPeer.mutex.Unlock()
	}

	if msg.MsgType == Reply {
		myPeer.mutex.Lock()

		if verbose {
			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp)
		}

		//aggiungo a replyProSet il msg
		myPeer.replySet.PushBack(msg)

		myPeer.mutex.Unlock()
	} else if msg.MsgType == Release {
		myPeer.mutex.Lock()

		if verbose {
			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp)
		}

		RemoveFirstElementMap(myPeer.ScalarMap)
		go checkAcks()
		myPeer.numRelease++
		if myPeer.numRelease == numSender {
			myPeer.StartTest <- true
		}
		myPeer.numMsgsTest++

		myPeer.mutex.Unlock()

	}

}

func checkAcks() {

	for !(myPeer.replySet.Len() == myPeer.PeerList.Len()-1 && len(myPeer.ScalarMap) > 0) {

		time.Sleep(time.Second * 5)
	}

	//prendo il primo mess nella mappa per vedere se è il mio, ossia guardo ID sender
	msg := GetFirstElementMap(myPeer.ScalarMap)

	if msg.Sender == myPeer.Username {
		myPeer.mutex.Lock()
		//il primo msg in lista è il mio, quindi posso accedere in CS

		if verbose {
			utilities.WriteInfosToFile("receives all peer reply messages successfully.", myPeer.LogPath, myPeer.Username)
		}

		//ho ricevuto tutti msg reply, ora entro in cs
		date := time.Now().Format(utilities.DateFormat)

		if verbose {
			utilities.WriteInfosToFile("enters the critical section at "+date+".", myPeer.LogPath, myPeer.Username)
		}

		time.Sleep(time.Minute / 2)
		date = time.Now().Format(utilities.DateFormat)

		if verbose {
			utilities.WriteInfosToFile("exits the critical section at "+date+".", myPeer.LogPath, myPeer.Username)
		}

		//lascio CS e mando msg release a tutti
		err := sendRelease()
		if err != nil {
			log.Fatalf("Error sending release msg %v", err)
		}

		myPeer.mutex.Unlock()
	}

}
