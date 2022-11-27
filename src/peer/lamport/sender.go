package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

var myPeer LamportPeer

func SendLamport(peer *LamportPeer) {

	if myPeer.Username == "" {
		fmt.Println("sto in SendLamport --- RA_PEER VUOTA")
		myPeer = *peer
		//peerCnt = MyRApeer.PeerList.Len()
	} else {
		fmt.Println("sto in SendLamport --- RA_PEER NON VUOTA")
	}

	myPeer.replySet.Init()
	// tale lista serve a mettere i msg di reply per poi controllare che sono arrivati tutti
	// TODO: invece che lista basta semplicemente un contatore?!?!

	myPeer.mutex.Lock()

	IncrementTS(&myPeer.Timestamp)
	fmt.Println("------------------ timestamp  ==", myPeer.Timestamp)

	date := time.Now().Format(utilities.DATE_FORMAT)

	msg := *NewRequest(myPeer.Username, date, myPeer.Timestamp)

	fmt.Println("IL MESSAGGIO E' ====", msg)
	//fmt.Println("ID MESSAGGIO E' ====", msg.MsgID)
	fmt.Println("MsgType MESSAGGIO E' ====", msg.MsgType)
	fmt.Println("Sender MESSAGGIO E' ====", msg.Sender)
	fmt.Println("Date MESSAGGIO E' ====", msg.Date)
	fmt.Println("timeStamp MESSAGGIO E' ====", msg.TS)
	sendRequest(msg)

	myPeer.Waiting = true

	myPeer.mutex.Unlock()

	WriteInfoToFile(myPeer.Username, myPeer.LogPath, "wait all peer reply messages.", false)

	<-myPeer.ChanAcquireLock

	WriteInfoToFile(myPeer.Username, myPeer.LogPath, " receive all peer reply messages successfully.", false)

	//ho ricevuto tutti msg reply, ora entro in cs
	fmt.Println("lista di msg in coda ==", myPeer.ScalarMap)
	fmt.Println("entro in CS")
	date = time.Now().Format(utilities.DATE_FORMAT)

	WriteInfoToFile(myPeer.Username, myPeer.LogPath, " enters the critical section at "+date+".", true)
	time.Sleep(time.Minute / 2)
	date = time.Now().Format(utilities.DATE_FORMAT)

	WriteInfoToFile(myPeer.Username, myPeer.LogPath, " exits the critical section at "+date+".", true)

	//lascio CS e mando msg release a tutti
	sendRelease()

}

func sendRequest(msg Message) error {

	WriteTSInfoToFile(myPeer.LogPath, myPeer.Username, myPeer.Timestamp, "lamport")

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
			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", msg, myPeer.Timestamp, "lamport")

			if err != nil {
				return err
			}

		}
	}
	//una volta inviato il msg, lo salvo nella coda locale del peer sender
	fmt.Println(" ------------------------------------------ STO QUA 2 ----------------------------")

	AppendHashMap(myPeer.ScalarMap, msg)
	fmt.Println(" ------------------------------------------ STO QUA 3 ----------------------------")

	/*
		for e := lista(); e != nil; e = e.Next() {
			item := e.Value.(Message)

		}
		fmt.Println("LISTA DEL PEER SENDER ==", lista)

	*/

	fmt.Println("MAPPA SENDER ====", myPeer.ScalarMap)

	return nil
}

func sendRelease() error {
	//incremento timestamp

	/*
		utilities.IncrementTS(&myPeer.Timestamp)
		utilities.WriteTSInfoToFile(myPeer.LogPath, myPeer.Username, myPeer.Timestamp, "lamport")
	*/
	date := time.Now().Format(utilities.DATE_FORMAT)
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

			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", releaseMsg, myPeer.Timestamp, "lamport")

			if err != nil {
				return err
			}
		}
	}

	//elimino primo msg da lista
	RemoveFirstElementMap(myPeer.ScalarMap)
	fmt.Println("ora la mappa ===", myPeer.ScalarMap)
	return nil
}
func sendReply(msg *Message) error {
	// mando ack al peer con id msg.receiver
	WriteTSInfoToFile(myPeer.LogPath, myPeer.Username, myPeer.Timestamp, "lamport")

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

			WriteMsgToFile(myPeer.LogPath, myPeer.Username, "send", *msg, myPeer.Timestamp, "lamport")

		}
	}
	return nil

}
