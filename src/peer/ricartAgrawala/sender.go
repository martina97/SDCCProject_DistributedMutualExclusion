package ricartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	MyRApeer RApeer
	peerCnt  int
)

func SendRicart(peer *RApeer) {

	if MyRApeer == (RApeer{}) {
		fmt.Println("sto in SendRicart --- RA_PEER VUOTA")
		MyRApeer = *peer
		peerCnt = MyRApeer.PeerList.Len()
	} else {
		fmt.Println("sto in SendRicart --- RA_PEER NON VUOTA")
	}

	MyRApeer.mutex.Lock()
	fmt.Println("sono in sendRicart!!!!! il peer ==", MyRApeer.ToString())
	for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
		fmt.Println("e ==", e)
	}

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

	fmt.Println("sono in sendRicard --- MyRApeer.Num ==== ", MyRApeer.Num)

	//	2. Num = Num+1; Last_Req = Num;
	utilities.IncrementTS(&MyRApeer.Num)
	MyRApeer.lastReq = MyRApeer.Num
	fmt.Println("sono in sendRicard --- MyRApeer.Num ==== ", MyRApeer.Num)
	fmt.Println(MyRApeer.ToString())

	//	3. for j=1 to N-1 send REQUEST to pj; --> INVIO MSG REQUEST AGLI ALTRI PEER
	date := time.Now().Format(utilities.DATE_FORMAT)
	msg := *utilities.NewRequest2(MyRApeer.Username, date, MyRApeer.lastReq)
	fmt.Println("IL MESSAGGIO E' ====", msg.ToString("send"))
	sendRequest(msg)

	utilities.WriteInfoToFile2(MyRApeer.Username, MyRApeer.LogPath, "wait all peer reply messages.", false)
	MyRApeer.mutex.Unlock()

	//4. Wait until #replies=N-1;
	<-MyRApeer.ChanAcquireLock

	utilities.WriteInfoToFile2(MyRApeer.Username, MyRApeer.LogPath, "receive all peer reply messages successfully.", false)
	//5. State = CS;
	MyRApeer.state = CS

	//6. CS
	fmt.Println("entro in CS")
	date = time.Now().Format(utilities.DATE_FORMAT)

	utilities.WriteInfoToFile2(MyRApeer.Username, MyRApeer.LogPath, " entered the critical section at "+date, true)
	time.Sleep(time.Minute / 2) //todo: invece che sleep mettere file condiviso
	date = time.Now().Format(utilities.DATE_FORMAT)
	utilities.WriteInfoToFile2(MyRApeer.Username, MyRApeer.LogPath, " exited the critical section at "+date, true)

	//7. ∀ r∈Q send REPLY to r
	fmt.Println("la lista dei msg in coda == ", MyRApeer.DeferSet)
	fmt.Println("la lista dei msg in coda MyRApeer.DeferSet.Front()== ", MyRApeer.DeferSet.Front())

	MyRApeer.mutex.Lock()
	MyRApeer.state = NCS
	for e := MyRApeer.DeferSet.Front(); e != nil; e = e.Next() {
		fmt.Println("msg ==", e.Value)
		//fmt.Println("msg ==", e)
		fmt.Println("msg ==", e.Value.(*utilities.Message))
		fmt.Println("-----")
		queueMsg := e.Value.(*utilities.Message)
		fmt.Println("queueMsg.sender = ", queueMsg.Sender)
		date := time.Now().Format(utilities.DATE_FORMAT)
		replyMsg := utilities.NewReply2(MyRApeer.Username, queueMsg.Sender, date, MyRApeer.Num)
		fmt.Println("il msg di REPLY ===", replyMsg.ToString("send"))

		for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
			dest := e.Value.(utilities.NodeInfo)
			if dest.Username == queueMsg.Sender {
				fmt.Println("invio msg reply a ---> ", queueMsg.Sender)
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

func sendRequest(msg utilities.Message) error {

	fmt.Println("sto in sendRequest")
	//scrivo sul log che ho aggiornato il TS
	//utilities.WriteTSInfoToFile(myID, MyRApeer.Num, algorithm)
	utilities.WriteTSInfoToFile2(MyRApeer.LogPath, MyRApeer.Username, MyRApeer.Num, "ricartAgrawala")

	fmt.Println("dopo WriteTSInfoToFile2")
	for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		fmt.Println("dest ==", dest)
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
			fmt.Println("dopo encode msg, msg == ", msg.ToString("send"))

			msg.Receiver = dest.Username
			fmt.Println("dopo encode msg, msg == ", msg.ToString("send"))

			//err = utilities.WriteMsgToFile(&myNode, "Send", msg, dest.ID, myNode.TimeStamp)
			//err = utilities.WriteMsgToFile2(MyRApeer.ID, "Send", msg, dest.ID, MyRApeer.Num, algorithm)
			err = utilities.WriteMsgToFile3(MyRApeer.LogPath, MyRApeer.Username, "send", msg, MyRApeer.Num, "ricartAgrawala")
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

func sendReply(msg *utilities.Message, receiver *utilities.NodeInfo) error {
	/*
		for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
			dest := e.Value.(utilities.NodeInfo)
			if dest.Username == msg.Receiver {

	*/
	fmt.Println("mando reply a ", msg.Receiver)
	fmt.Println("receiver = ", receiver.Username)
	//open connection whit peer
	peerConn := receiver.Address + ":" + receiver.Port
	conn, err := net.Dial("tcp", peerConn)
	defer conn.Close()
	if err != nil {
		log.Println("Send response error on Dial")
	}
	enc := gob.NewEncoder(conn)
	enc.Encode(msg)

	f, err := os.OpenFile(MyRApeer.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)
	_, err = f.WriteString("[" + date + "] : " + MyRApeer.Username + " send" + msg.ToString("send") + " to " + receiver.Username)
	_, err = f.WriteString("\n")
	//err = utilities.WriteMsgToFile3(MyRApeer.LogPath, MyRApeer.Username, "Send", msg, MyRApeer.Num, "ricartAgrawala")

	err = f.Sync()
	if err != nil {
		return err
	}
	/*
			}
		}

	*/
	return nil
}
