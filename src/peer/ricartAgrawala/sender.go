package ricartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	MyRApeer RApeer
	peerCnt  int
)

func SendRicart(peer *RApeer) {

	if MyRApeer == (RApeer{}) {
		MyRApeer = *peer
		peerCnt = MyRApeer.PeerList.Len()
	}

	MyRApeer.mutex.Lock()

	//inizializzo le variabili che mi servono
	MyRApeer.DeferSet.Init()
	//fmt.Println("MyRApeer.DeferSet.Init()")
	MyRApeer.replySet.Init()
	//fmt.Println("MyRApeer.replySet.Init()")
	MyRApeer.replies = 0
	//fmt.Println("MyRApeer.replies = 0")

	/*
		1. State = Requesting;
		2. Num = Num+1; Last_Req = Num;
		3. for j=1 to N-1 send REQUEST to pj;
		4. Wait until #replies=N-1;
		5. State = CS;
		6. CS
		7. ∀ r∈Q send REPLY to r
		8. Q=∅; State=NCS; #replies=0;
	*/

	// 1. State = Requesting;
	MyRApeer.state = Requesting

	//	2. Num = Num+1; Last_Req = Num;
	lamport.IncrementTS(&MyRApeer.Num)
	MyRApeer.lastReq = MyRApeer.Num

	//	3. for j=1 to N-1 send REQUEST to pj; --> INVIO MSG REQUEST AGLI ALTRI PEER
	date := time.Now().Format(utilities.DATE_FORMAT)
	msg := *lamport.NewRequest(MyRApeer.Username, date, MyRApeer.lastReq)
	sendRequest(msg)
	MyRApeer.mutex.Unlock()

	utilities.WriteInfosToFile("waits all peer reply messages.", MyRApeer.LogPath, MyRApeer.Username)

	//4. Wait until #replies=N-1;
	<-MyRApeer.ChanAcquireLock

	utilities.WriteInfosToFile("receives all peer reply messages successfully.", MyRApeer.LogPath, MyRApeer.Username)

	//5. State = CS;
	MyRApeer.state = CS

	//6. CS
	date = time.Now().Format(utilities.DATE_FORMAT)

	utilities.WriteInfosToFile(" enters the critical section at "+date+".", MyRApeer.LogPath, MyRApeer.Username)

	time.Sleep(time.Minute / 2) //todo: invece che sleep mettere file condiviso
	date = time.Now().Format(utilities.DATE_FORMAT)
	utilities.WriteInfosToFile(" exits the critical section at "+date+".", MyRApeer.LogPath, MyRApeer.Username)

	//7. ∀ r∈Q send REPLY to r

	MyRApeer.mutex.Lock()
	MyRApeer.state = NCS
	//todo: se DeferSet vuota?
	for e := MyRApeer.DeferSet.Front(); e != nil; e = e.Next() {

		queueMsg := e.Value.(*lamport.Message)
		date := time.Now().Format(utilities.DATE_FORMAT)
		replyMsg := lamport.NewReply(MyRApeer.Username, queueMsg.Sender, date, MyRApeer.Num)

		for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
			dest := e.Value.(utilities.NodeInfo)
			if dest.Username == queueMsg.Sender {
				// invio msg reply a queueMsg.Sender
				err := sendReply(replyMsg, &dest)
				if err != nil {
					log.Fatalf("error sending ack %v", err)
				}
			}
		}

	}
	MyRApeer.mutex.Unlock()
	//8. Q=∅; State=NCS; #replies=0;

}

func sendRequest(msg lamport.Message) error {

	utilities.WriteTSInfoToFile(MyRApeer.LogPath, MyRApeer.Username, strconv.Itoa(int(MyRApeer.Num)))

	for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		if dest.Type == utilities.Peer && dest.ID != MyRApeer.ID { //non voglio mandarlo a me stesso

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

			lamport.WriteMsgToFile(MyRApeer.LogPath, MyRApeer.Username, "send", msg, MyRApeer.Num)

		}
	}

	return nil
}

func sendReply(msg *lamport.Message, receiver *utilities.NodeInfo) error {
	//mando reply a msg.Receiver
	peerConn := receiver.Address + ":" + receiver.Port
	conn, err := net.Dial("tcp", peerConn)
	defer conn.Close()
	if err != nil {
		log.Println("Send response error on Dial")
	}
	enc := gob.NewEncoder(conn)
	enc.Encode(msg)

	lamport.WriteMsgToFile(MyRApeer.LogPath, MyRApeer.Username, "send", *msg, MyRApeer.Num)
	
	return nil
}
