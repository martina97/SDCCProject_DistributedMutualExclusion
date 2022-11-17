package RicartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	MyRApeer RApeer
)

func SendRicart(peer RApeer) {
	MyRApeer = peer
	fmt.Println("sono in sendRicart!!!!! il peer ==", MyRApeer.ToString())
	for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
		fmt.Println("e ==", e)
	}

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

	mu := MyRApeer.GetMutex()
	mu.Lock()
	fmt.Println("sono in sendRicard --- MyRApeer.Num ==== ", MyRApeer.Num)

	//	2. Num = Num+1; Last_Req = Num;
	utilities.IncrementTS(&MyRApeer.Num)
	MyRApeer.lastReq = MyRApeer.Num
	fmt.Println("sono in sendRicard --- MyRApeer.Num ==== ", MyRApeer.Num)
	fmt.Println(MyRApeer.ToString())

	//	3. for j=1 to N-1 send REQUEST to pj; --> INVIO MSG REQUEST AGLI ALTRI PEER
	date := time.Now().Format(utilities.DATE_FORMAT)
	msg := *utilities.NewRequest2(MyRApeer.Username, date, MyRApeer.lastReq)
	fmt.Println("IL MESSAGGIO E' ====", msg)
	sendRicartAgrawalaRequest(msg)

	mu.Unlock()
	utilities.WriteInfoToFile2(MyRApeer.Username, MyRApeer.LogPath, " wait all peer reply messages.", false)

	//4. Wait until #replies=N-1;
	<-MyRApeer.ChanAcquireLock

	utilities.WriteInfoToFile2(MyRApeer.Username, MyRApeer.LogPath, " receive all peer reply messages successfully.", false)
	//5. State = CS;

	//6. CS

	//7. ∀ r∈Q send REPLY to r

	//8. Q=∅; State=NCS; #replies=0;

}

func sendRicartAgrawalaRequest(msg utilities.Message) error {

	fmt.Println("sto in sendRicartAgrawalaRequest")
	//scrivo sul log che ho aggiornato il TS
	//utilities.WriteTSInfoToFile(myID, MyRApeer.Num, algorithm)
	utilities.WriteTSInfoToFile2(MyRApeer.LogPath, MyRApeer.Username, MyRApeer.Num, "RicartAgrawala")

	fmt.Println("dopo WriteTSInfoToFile2")
	for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		//only peer are destination of msgs
		if dest.Type == utilities.Peer && dest.ID != MyRApeer.ID { //non voglio mandarlo a me stesso

			//open connection whit peer
			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			msg.Receiver = dest.Username

			//err = utilities.WriteMsgToFile(&myNode, "Send", msg, dest.ID, myNode.TimeStamp)
			//err = utilities.WriteMsgToFile2(MyRApeer.ID, "Send", msg, dest.ID, MyRApeer.Num, algorithm)
			err = utilities.WriteMsgToFile3(MyRApeer.LogPath, MyRApeer.Username, "Send", msg, MyRApeer.Num, "RicartAgrawala")
			if err != nil {
				return err
			}

		}
	}
	/*
		//una volta inviato il msg, lo salvo nella coda locale del peer sender
		fmt.Println(" ------------------------------------------ STO QUA 2 ----------------------------")

		utilities.AppendHashMap2(myNode.ScalarMap, msg)
		fmt.Println(" ------------------------------------------ STO QUA 3 ----------------------------")

		fmt.Println("MAPPA SENDER ====", myNode.ScalarMap)

	*/

	return nil
}
