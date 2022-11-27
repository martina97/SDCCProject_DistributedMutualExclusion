package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"time"
)

var myPeer LamportPeer

func SendLamport(peer *LamportPeer) {

	if myPeer.Username == "" {
		myPeer = *peer
	}

	myPeer.replySet.Init()
	// tale lista serve a mettere i msg di reply per poi controllare che sono arrivati tutti

	myPeer.mutex.Lock()

	IncrementTS(&myPeer.Timestamp)

	date := time.Now().Format(utilities.DateFormat)

	msg := *NewRequest(myPeer.Username, date, myPeer.Timestamp)

	sendRequest(msg)

	myPeer.Waiting = true

	myPeer.mutex.Unlock()

	utilities.WriteInfosToFile("waits all peer reply messages.", myPeer.LogPath, myPeer.Username)

	go checkAcks()
	/*
		<-myPeer.ChanAcquireLock

		utilities.WriteInfosToFile("receives all peer reply messages successfully.", myPeer.LogPath, myPeer.Username)

		//ho ricevuto tutti msg reply, ora entro in cs
		date = time.Now().Format(utilities.DateFormat)

		utilities.WriteInfosToFile("enters the critical section at "+date+".", myPeer.LogPath, myPeer.Username)

		time.Sleep(time.Minute / 2)
		date = time.Now().Format(utilities.DateFormat)

		utilities.WriteInfosToFile("exits the critical section at "+date+".", myPeer.LogPath, myPeer.Username)

		//lascio CS e mando msg release a tutti
		sendRelease()

	*/

}

func sendRequest(msg Message) error {

	utilities.WriteTSInfoToFile(myPeer.LogPath, myPeer.Username, strconv.Itoa(int(myPeer.Timestamp)))

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		//only peer are destination of msgs
		if dest.Type == utilities.Peer && dest.ID != myPeer.ID { //non voglio mandarlo a me stesso

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

			//r = utilities.WriteMsgToFile(&myPeer, "Send", msg, dest.ID, myPeer.timestamp)
			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", msg, myPeer.Timestamp)

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

	date := time.Now().Format(utilities.DateFormat)
	releaseMsg := *NewRelease(myPeer.Username, date, myPeer.Timestamp)

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		//only peer are destination of msgs
		if dest.Type == utilities.Peer && dest.ID != myPeer.ID { //non voglio mandarlo a me stesso

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

			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", releaseMsg, myPeer.Timestamp)

			if err != nil {
				return err
			}
		}
	}

	//elimino primo msg da lista
	RemoveFirstElementMap(myPeer.ScalarMap)
	return nil
}
func sendReply(msg *Message) error {
	// mando ack al peer con id msg.receiver
	utilities.WriteTSInfoToFile(myPeer.LogPath, myPeer.Username, strconv.Itoa(int(myPeer.Timestamp)))

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		if dest.Username == msg.Receiver {
			//open connection with peer
			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", *msg, myPeer.Timestamp)

		}
	}
	return nil

}
