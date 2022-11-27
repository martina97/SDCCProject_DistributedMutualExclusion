package ricartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
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

	if MyRApeer == (RApeer{}) {
		MyRApeer = *peer
		peerCnt = MyRApeer.PeerList.Len()
	}

	defer conn.Close()
	msg := new(lamport.Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)

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
		MyRApeer.mutex.Lock()
		lamport.UpdateTS(&MyRApeer.Num, &msg.TS)

		lamport.WriteMsgToFile(MyRApeer.LogPath, MyRApeer.Username, "receive", *msg, MyRApeer.Num)

		if checkConditions(msg) { //se è true --> inserisco msg in coda
			MyRApeer.DeferSet.PushBack(msg)
		} else { //se è false --> invio REPLY al peer che ha inviato msg REQUEST
			date := time.Now().Format(utilities.DATE_FORMAT)
			replyMsg := lamport.NewReply(MyRApeer.Username, msg.Sender, date, MyRApeer.Num)

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
		MyRApeer.mutex.Lock()
		MyRApeer.replies++
		lamport.WriteMsgToFile(MyRApeer.LogPath, MyRApeer.Username, "receive", *msg, MyRApeer.Num)
		
		if MyRApeer.replies == peerCnt-1 {
			//ho ricevuto tutti i msg di reply
			MyRApeer.ChanAcquireLock <- true
		}
		MyRApeer.mutex.Unlock()

	}

	return nil
}

func checkConditions(msg *lamport.Message) bool {

	if (MyRApeer.state == CS) || (MyRApeer.state == Requesting && checkTS(msg)) {
		//non invio reply e metto msg in coda
		return true
	}
	//invio reply
	return false

}

func checkTS(msg *lamport.Message) bool {
	// true se {Last_Req, i} < {t, j})
	if (MyRApeer.lastReq <= msg.TS) && (MyRApeer.Username < msg.Sender) {
		// la condizione è true --> non invio reply e metto msg in coda
		return true
	}
	return false

}
