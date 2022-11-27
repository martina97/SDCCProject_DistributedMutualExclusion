package ricartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

/*
func Message_handler() {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.Client_port))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go HandleConnection(connection)
		//go handleConnectionCentralized(connection)
	}
}

*/

//Save message
func HandleConnection(conn net.Conn, peer *RApeer) error {

	fmt.Println("sto in handleConnection dentro ricartAgrawala package")
	if MyRApeer == (RApeer{}) {
		fmt.Println("RA_PEER VUOTA")
		MyRApeer = *peer
		peerCnt = MyRApeer.PeerList.Len()
	} else {
		fmt.Println("RA_PEER NON VUOTA")
	}

	//devo vedere se già ho inizializzato RApeer (se non ho mai inviato un msg, non l'ho inizializzato)
	fmt.Println("MyRApeer == ", MyRApeer.ToString())

	defer conn.Close()
	msg := new(lamport.Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)
	fmt.Println("il msg == ", msg.ToString("receive"))

	//mutex := MyRApeer.GetMutex()
	if msg.MsgType == lamport.Request {

		/*
			Upon receipt REQUEST(t) from pj
			1. if State=CS or (State=Requesting and {Last_Req, i} < {t, j})
				then insert {t, j} in Q
			3. else
				send REPLY to pj
			4. Num = max(t, Num)
		*/
		fmt.Println("MESS REQUEST !!!!!! ")
		MyRApeer.mutex.Lock()
		lamport.UpdateTS(&MyRApeer.Num, &msg.TS)

		utilities.WriteMsgToFile(MyRApeer.LogPath, MyRApeer.Username, "receive", *msg, MyRApeer.Num, "ricartAgrawala")

		if checkConditions(msg) { //se è true --> inserisco msg in coda
			MyRApeer.DeferSet.PushBack(msg)
		} else { //se è false --> invio REPLY al peer che ha inviato msg REQUEST
			date := time.Now().Format(utilities.DATE_FORMAT)
			replyMsg := lamport.NewReply(MyRApeer.Username, msg.Sender, date, MyRApeer.Num)
			fmt.Println("il msg di REPLY ===", replyMsg.ToString("send"))

			//devo inviare reply al replyMsg.receiver
			for e := MyRApeer.PeerList.Front(); e != nil; e = e.Next() {
				dest := e.Value.(utilities.NodeInfo)
				if dest.Username == replyMsg.Receiver {
					err := sendReply(replyMsg, &dest)
					if err != nil {
						log.Fatalf("error sending ack %v", err)
					}
				}
			}

		}
		MyRApeer.mutex.Unlock()
	}
	if msg.MsgType == lamport.Reply {
		/*
			Upon receipt REPLY from pj
			1. #replies = #replies+1
		*/
		fmt.Println("MESS REPLY !!!!!! ")
		MyRApeer.mutex.Lock()
		//MyRApeer.replies =MyRApeer.replies + 1
		MyRApeer.replies++
		fmt.Println("replies = ", MyRApeer.replies)
		utilities.WriteMsgToFile(MyRApeer.LogPath, MyRApeer.Username, "receive", *msg, MyRApeer.Num, "ricartAgrawala")
		fmt.Println("peerCnt = ", peerCnt)

		if MyRApeer.replies == peerCnt-1 {
			fmt.Println("ho ricevuto tutti i msg di reply")
			MyRApeer.ChanAcquireLock <- true
		}
		MyRApeer.mutex.Unlock()

	}

	return nil
}

func checkConditions(msg *lamport.Message) bool {

	if (MyRApeer.state == CS) || (MyRApeer.state == Requesting && checkTS(msg)) {
		fmt.Println("sto in checkConditions -->  non invio reply e metto msg in coda")
		return true
	}
	fmt.Println("invio reply!!!!!!")
	return false

}

func checkTS(msg *lamport.Message) bool {
	// true se {Last_Req, i} < {t, j})
	if (MyRApeer.lastReq <= msg.TS) && (MyRApeer.Username < msg.Sender) {
		fmt.Println("sto in checkTS e la condizione e' true --> non invio reply e metto msg in coda")
		return true
	}
	return false

}
