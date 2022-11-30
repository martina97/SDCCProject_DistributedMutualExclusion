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
	MyRApeer.replies = 0

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
	lamport.IncrementSC(&MyRApeer.Num)
	MyRApeer.lastReq = MyRApeer.Num

	//	3. for j=1 to N-1 send REQUEST to pj; --> INVIO MSG REQUEST AGLI ALTRI PEER
	date := time.Now().Format(utilities.DateFormat)
	msg := *lamport.NewRequest(MyRApeer.Username, date, MyRApeer.lastReq)
	sendRequest(msg)
	MyRApeer.mutex.Unlock()

	utilities.WriteInfosToFile("waits all peer reply messages.", MyRApeer.LogPath, MyRApeer.Username)

	//4. Wait until #replies=N-1;

	go checkAcks()

}

func sendRequest(msg lamport.Message) error {

	utilities.WriteTSInfoToFile(MyRApeer.LogPath, MyRApeer.Username, strconv.Itoa(int(MyRApeer.Num)))

	for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		if dest.Type == utilities.Peer && dest.ID != MyRApeer.ID { //non voglio mandarlo a me stesso
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

			lamport.WriteMsgToFile(MyRApeer.LogPath, MyRApeer.Username, "send", msg, MyRApeer.Num)

		}
	}

	return nil
}

func sendReply(msg *lamport.Message, receiver *utilities.NodeInfo) error {
	//mando reply a msg.Receiver
	utilities.SleepRandInt()

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
