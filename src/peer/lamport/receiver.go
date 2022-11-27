package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
)

func HandleConnection(conn net.Conn, peer *LamportPeer) {
	fmt.Println("sto in handleConnection dentro lamport package")
	if myPeer.Username == "" {
		fmt.Println("LamportPeer VUOTA")
		myPeer = *peer
		//peerCnt = MyRApeer.PeerList.Len()
	} else {
		fmt.Println("LamportPeer NON VUOTA")
	}

	// read msg and save on file
	defer conn.Close()
	msg := new(Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)
	fmt.Println("il msg == ", msg.ToString("receive"))

	//ogni volta che ricevo un msg devo aggiornare TS
	//aggiorno timestamp
	//ogni peer ha il suo clock scalare, e' var globale come myNode e myID

	//time.Sleep(time.Minute / 2) //PRIMA DI AUMENTARE TS METTO SLEEP COSI PROVO A INVIARE 2 REQ INSIEME E VEDO CHE SUCCEDE

	//mutex := lock.GetMutex()

	if msg.MsgType == Request {
		UpdateTS(&myPeer.Timestamp, &msg.TS)
		/*
			quando ricevo una richiesta da un processo devo decidere se mandare ACK al processo oppure se voglio entrare in CS
		*/
		fmt.Println("MESS REQUEST !!!!!! ")
		fmt.Println("TIMESTAMP QUANDO RICEVO REQUEST ===", myPeer.Timestamp) //ho gia aggiornato il TS!!
		//fmt.Println("------------------------------------------------------------- DOPO RICEVUTO REQUEST --- > timestamp  ==", timeStamp)
		myPeer.mutex.Lock()
		utilities.WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp, "lamport")

		//metto msg in mappa
		AppendHashMap(myPeer.ScalarMap, *msg)

		//QUA DEVO DECIDERE SE MANDARE ACK O REQUEST (msg REPLY O REQUEST)

		//scelta := openMenuRequest()
		//fmt.Println("HO SCELTO", scelta)

		//mando msg reply
		//date := time.Now().Format("17:06:04")
		//prima di mandare reply aggiorno di nuovo il TS !!
		//utilities.IncrementTS(&myPeer.Timestamp)

		fmt.Println("------------------------------------------------------------- DOPO INVIATO REPLY --- > timestamp  ==", myPeer.Timestamp)
		date := time.Now().Format(utilities.DATE_FORMAT)
		replyMsg := NewReply(myPeer.Username, msg.Sender, date, myPeer.Timestamp)
		sendReply(replyMsg)
		myPeer.mutex.Unlock()
	}

	if msg.MsgType == Reply {
		fmt.Println("------------------------------------------------------------- DOPO RICEVUTO REPLY --- > timestamp  ==", myPeer.Timestamp)
		myPeer.mutex.Lock()
		fmt.Println("TIMESTAMP QUANDO RICEVO Reply ===", myPeer.Timestamp)

		//utilities.WriteMsgToFile(&myNode, "Receive", *msg, 0, myNode.TimeStamp)
		utilities.WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp, "lamport")

		//utilities.WriteTSInfoToFile(myID, timeStamp)

		/*
			f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}


				//save new address on file
				date := time.Now().Format("15:04:05.000")
				_, err = f.WriteString("[" + date + "] : Receive" + msg.ToString())
				_, err = f.WriteString("\n")
				err = f.Sync()
				if err != nil {
					return err
				}
		*/

		//aggiungo a replyProSet il msg
		myPeer.replySet.PushBack(msg)
		//check ack
		checkAcks() //controllo se ho ricevuto 2 msg reply, se si posso entrare in CS prendendo 1 elem nella lista
		// e controllando che id sia il mio, se e' il mio entro altrimenti no
		//todo: sez critica?!?!??!
		myPeer.mutex.Unlock()
	} else if msg.MsgType == Release {
		fmt.Println(" RICEVO RELEASE !! ")
		myPeer.mutex.Lock()

		//utilities.WriteMsgToFile(&myNode, "Receive", *msg, 0, myNode.TimeStamp)
		utilities.WriteMsgToFile(myPeer.LogPath, myPeer.Username, "receive", *msg, myPeer.Timestamp, "lamport")
		fmt.Println("ho scritto su file")
		//utilities.WriteTSInfoToFile(myID, timeStamp)

		RemoveFirstElementMap(myPeer.ScalarMap)
		fmt.Println("---------------------------------	DOPO AVER RICEVUTO RELEASE mappa ===", myPeer.ScalarMap)
		checkAcks()
		myPeer.numRelease++
		if myPeer.numRelease == numSender {
			myPeer.StartTest <- true
		}
		myPeer.mutex.Unlock()

	}

	fmt.Println("msg ricevuti -----", myPeer.ScalarMap)

}

func checkAcks() {

	//todo: quando azzero lista ReplyProSet ?????
	//date := time.Now().Format("15:04:05.000")
	fmt.Println("SONO IN checkAcks")

	fmt.Println("process.GetReplyProSet().Len() ==== ", myPeer.replySet.Len())
	fmt.Println("peers.Len()-1 ==== ", myPeer.PeerList.Len()-1)
	fmt.Println("len(scalarMap) ==== ", len(myPeer.ScalarMap))

	if myPeer.replySet.Len() == myPeer.PeerList.Len()-1 && len(myPeer.ScalarMap) > 0 {
		fmt.Println("HO RICEVUTO 2 MSG REPLY")

		//prendo il primo mess nella mappa per vedere se è il mio, ossia guardo ID sender
		msg := GetFirstElementMap(myPeer.ScalarMap)
		fmt.Println("MSG IN CHECK ACKS ===", msg)

		if msg.Sender == myPeer.Username {
			//il primo msg in lista è il mio, quindi posso accedere in CS
			myPeer.Waiting = false
			myPeer.ChanAcquireLock <- true
		}

	}
}
