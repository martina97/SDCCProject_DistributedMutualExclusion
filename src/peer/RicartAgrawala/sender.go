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
	myRApeer RApeer
)

func SendRicart(peer *RApeer) {
	myRApeer = *peer
	fmt.Println("sono in sendRicart!!!!! il peer ==", myRApeer.ToString())
	for e := myRApeer.PeerList.Front(); e != nil; e = e.Next() {
		fmt.Println("e == ", e)

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
	myRApeer.state = Requesting

	mu := myRApeer.GetMutex()
	mu.Lock()
	fmt.Println("sono in sendRicard --- myRApeer.Num ==== ", myRApeer.Num)

	//	2. Num = Num+1; Last_Req = Num;
	utilities.IncrementTS(&myRApeer.Num)
	myRApeer.lastReq = myRApeer.Num
	fmt.Println("sono in sendRicard --- myRApeer.Num ==== ", myRApeer.Num)
	fmt.Println(myRApeer.ToString())

	//	3. for j=1 to N-1 send REQUEST to pj; --> INVIO MSG REQUEST AGLI ALTRI PEER
	date := time.Now().Format(utilities.DATE_FORMAT)
	msg := *utilities.NewRequest2(myRApeer.Username, date, myRApeer.lastReq)
	fmt.Println("IL MESSAGGIO E' ====", msg)
	sendRicartAgrawalaRequest(msg)

	mu.Unlock()
	utilities.WriteInfoToFile2(myRApeer.Username, myRApeer.LogPath, " wait all peer reply messages.", false)

	//4. Wait until #replies=N-1;
	<-myRApeer.ChanAcquireLock

	utilities.WriteInfoToFile2(myRApeer.Username, myRApeer.LogPath, " receive all peer reply messages successfully.", false)
	//5. State = CS;

	//6. CS

	//7. ∀ r∈Q send REPLY to r

	//8. Q=∅; State=NCS; #replies=0;

}

func sendRicartAgrawalaRequest(msg utilities.Message) error {

	fmt.Println("sto in sendRicartAgrawalaRequest")
	//scrivo sul log che ho aggiornato il TS
	//utilities.WriteTSInfoToFile(myID, myRApeer.Num, algorithm)
	utilities.WriteTSInfoToFile2(myRApeer.LogPath, myRApeer.Username, myRApeer.Num, "RicartAgrawala")

	for e := myRApeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		//only peer are destination of msgs
		if dest.Type == utilities.Peer && dest.ID != myRApeer.ID { //non voglio mandarlo a me stesso

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
			//err = utilities.WriteMsgToFile2(myRApeer.ID, "Send", msg, dest.ID, myRApeer.Num, algorithm)
			err = utilities.WriteMsgToFile3(myRApeer.LogPath, myRApeer.Username, "Send", msg, myRApeer.Num, "RicartAgrawala")
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
