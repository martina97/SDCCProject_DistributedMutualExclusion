package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

var myPeer LamportPeer

func SendLamport(peer *LamportPeer) {

	/*
		flag.BoolVar(&verbose, "v", utilities.Verbose, "use this flag to get verbose info on messages")
		flag.Parse()

		if verbose {
			fmt.Println("VERBOSE FLAG ON")
		}

	*/

	if myPeer.Username == "" {
		myPeer = *peer
	}

	myPeer.replySet.Init()
	// tale lista serve a mettere i msg di reply per poi controllare che sono arrivati tutti

	myPeer.mutex.Lock()

	IncrementSC(&myPeer.Timestamp)

	date := time.Now().Format(utilities.DateFormat)

	msg := *NewRequest(myPeer.Username, date, myPeer.Timestamp)
	utilities.SleepRandInt()

	sendRequest(msg)

	myPeer.Waiting = true

	myPeer.mutex.Unlock()

	if verbose {
		utilities.WriteInfosToFile("waits all peer reply messages.", myPeer.LogPath, myPeer.Username)
	}

	go checkAcks()
	if utilities.Test {
		go checkStartTests()
	}

}

func checkStartTests() {
	fmt.Println("sto in checkStartTests")
	for !(myPeer.numMsgsTest == numSender) {
		fmt.Println("sto in checkStartTests dentro for")

		time.Sleep(time.Second * 5)
	}
	fmt.Println("sto in checkStartTests fuori for")

	myPeer.ChanStartTest <- true
	fmt.Println("chan true!!!!")

}

func sendRequest(msg Message) error {

	if verbose {
		utilities.WriteTSInfoToFile(myPeer.LogPath, myPeer.Username, strconv.Itoa(int(myPeer.Timestamp)))
	}

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		if dest.Type == utilities.Peer && dest.ID != myPeer.ID { //non voglio mandarlo a me stesso

			utilities.SleepRandInt()
			//open connection with peer
			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			msg.Receiver = dest.Username

			if verbose {
				WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", msg, myPeer.Timestamp)
			}

			if err != nil {
				return err
			}

		}
	}
	//una volta inviato il msg, lo salvo nella coda locale del peer sender

	AppendHashMap(myPeer.ScalarMap, msg)

	return nil
}

func sendRelease() error {
	//incremento timestamp
	fmt.Println("mando release")

	date := time.Now().Format(utilities.DateFormat)
	releaseMsg := *NewRelease(myPeer.Username, date, myPeer.Timestamp)

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		//only peer are destination of msgs
		if dest.Type == utilities.Peer && dest.ID != myPeer.ID { //non voglio mandarlo a me stesso
			utilities.SleepRandInt()
			//open connection with peer
			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(releaseMsg)

			releaseMsg.Receiver = dest.Username

			if verbose {
				WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", releaseMsg, myPeer.Timestamp)
			}

			if err != nil {
				return err
			}
		}
	}

	//elimino primo msg da lista
	RemoveFirstElementMap(myPeer.ScalarMap)
	myPeer.numMsgsTest++
	return nil
}
func sendReply(msg *Message) error {
	// mando ack al peer con id msg.receiver
	if verbose {
		utilities.WriteTSInfoToFile(myPeer.LogPath, myPeer.Username, strconv.Itoa(int(myPeer.Timestamp)))
	}

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		if dest.Username == msg.Receiver {
			utilities.SleepRandInt()
			//open connection with peer
			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			if verbose {
				WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", *msg, myPeer.Timestamp)
			}

		}
	}
	return nil

}
