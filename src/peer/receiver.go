package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	msgScaFile *os.File
)

func message_handler_centr() {
	//listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.Client_port))
	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		log.Fatal("tcp.Lister fail")
	}
	defer listener.Close()

	//open file for save msg
	open_files()
	defer close_files()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		//go handleConnection(connection)
		go handleConnectionCentralized(connection)
	}
}

func message_handler() {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.Client_port))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	//open file for save msg
	open_files()
	defer close_files()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		switch algorithm {
		case "RicartAgrawala":

		}
		go handleConnection(connection)
		//go handleConnectionCentralized(connection)
	}
}

func open_files() {
	var err error
	msgScaFile, err = os.OpenFile(utilities.Peer_msg_sca_file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Impossible to open file")
	}
}

func close_files() {

	err := msgScaFile.Close()
	if err != nil {
		log.Fatal(err)
	}

}

//Save message
func handleConnection(conn net.Conn) error {
	// read msg and save on file
	defer conn.Close()
	msg := new(utilities.Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)

	//ogni volta che ricevo un msg devo aggiornare TS
	//aggiorno timestamp
	tmp := msg.SeqNum
	//ogni peer ha il suo clock scalare, e' var globale come myNode e myID

	time.Sleep(time.Minute / 2) //PRIMA DI AUMENTARE TS METTO SLEEP COSI PROVO A INVIARE 2 REQ INSIEME E VEDO CHE SUCCEDE

	utilities.UpdateTS(&myNode.TimeStamp, &msg.TS)

	//mutex := lock.GetMutex()
	mutex := myNode.GetMutex()
	if msg.MsgType == utilities.Request {
		/*
			quando ricevo una richiesta da un processo devo decidere se mandare ACK al processo oppure se voglio entrare in CS
		*/
		fmt.Println("MESS REQUEST !!!!!! ")
		fmt.Println("TIMESTAMP QUANDO RICEVO REQUEST ===", myNode.TimeStamp) //ho gia aggiornato il TS!!
		//fmt.Println("------------------------------------------------------------- DOPO RICEVUTO REQUEST --- > timestamp  ==", timeStamp)

		mutex.Lock()
		utilities.WriteMsgToFile(&myNode, "Receive", *msg, 0, myNode.TimeStamp)
		//utilities.WriteTSInfoToFile(myID, timeStamp)
		/*
			f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			//save new address on file
			date := time.Now().Format("15:04:05.000")
			_, err = f.WriteString("[" + date + "] : Receive " + msg.MessageToString())
			_, err = f.WriteString("\n")
			err = f.Sync()
			if err != nil {
				return err
			}

		*/

		//metto msg in mappa
		utilities.AppendHashMap2(myNode.ScalarMap, *msg)

		//QUA DEVO DECIDERE SE MANDARE ACK O REQUEST (msg REPLY O REQUEST)

		//scelta := openMenuRequest()
		//fmt.Println("HO SCELTO", scelta)

		//mando msg reply
		//date := time.Now().Format("17:06:04")
		//prima di mandare reply aggiorno di nuovo il TS !!
		utilities.IncrementTS(&myNode.TimeStamp)

		fmt.Println("------------------------------------------------------------- DOPO INVIATO REPLY --- > timestamp  ==", myNode.TimeStamp)
		date := time.Now().Format(utilities.DATE_FORMAT)
		replyMsg := utilities.NewReply(tmp, myNode.Username, msg.Sender, date, myNode.TimeStamp)
		sendAck(replyMsg)
		mutex.Unlock()
	}

	if msg.MsgType == utilities.Reply {
		fmt.Println("------------------------------------------------------------- DOPO RICEVUTO REPLY --- > timestamp  ==", myNode.TimeStamp)
		mutex.Lock()
		fmt.Println("TIMESTAMP QUANDO RICEVO Reply ===", myNode.TimeStamp)

		utilities.WriteMsgToFile(&myNode, "Receive", *msg, 0, myNode.TimeStamp)
		//utilities.WriteTSInfoToFile(myID, timeStamp)

		/*
			f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}


				//save new address on file
				date := time.Now().Format("15:04:05.000")
				_, err = f.WriteString("[" + date + "] : Receive" + msg.MessageToString())
				_, err = f.WriteString("\n")
				err = f.Sync()
				if err != nil {
					return err
				}
		*/

		//aggiungo a replyProSet il msg
		myNode.GetReplyProSet().PushBack(msg)
		//check ack
		checkAcks(&myNode) //controllo se ho ricevuto 2 msg reply, se si posso entrare in CS prendendo 1 elem nella lista
		// e controllando che id sia il mio, se e' il mio entro altrimenti no
		//todo: sez critica?!?!??!
		mutex.Unlock()
	} else if msg.MsgType == utilities.Release {
		mutex.Lock()

		/*
			date := time.Now().Format("15:04:05.000")
			f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			//save new msg on file
			_, err = f.WriteString("[" + date + "] : Receive " + msg.MessageToString())
			_, err = f.WriteString("\n")
			err = f.Sync()
				if err != nil {
					return err
				}

		*/
		utilities.WriteMsgToFile(&myNode, "Receive", *msg, 0, myNode.TimeStamp)

		//utilities.WriteTSInfoToFile(myID, timeStamp)

		utilities.RemoveFirstElementMap(myNode.ScalarMap)
		fmt.Println("---------------------------------	DOPO AVER RICEVUTO RELEASE mappa ===", myNode.ScalarMap)
		checkAcks(&myNode)
		mutex.Unlock()

	}

	fmt.Println("msg ricevuti -----", myNode.ScalarMap)
	/*
		for key, element := range scalarMap {
			fmt.Println("Key:", key, "=>", "Element:", element)
		}

	*/
	//fmt.Println("PRIMO MESS ==", utilities.GetFirstElementMap(scalarMap))

	return nil
}

//Save message
func handleConnectionCentralized(conn net.Conn) error {
	// read msg and save on file
	defer conn.Close()
	msg := new(utilities.CentralizedMessage)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)
	fmt.Println("sto in handleConnectionCentralized")
	fmt.Println("msg == ", msg)

	//time.Sleep(time.Minute / 2) //PRIMA DI AUMENTARE TS METTO SLEEP COSI PROVO A INVIARE 2 REQ INSIEME E VEDO CHE SUCCEDE

	//mutex := lock.GetMutex()
	//	mutex := myNode.GetMutex()
	if msg.MsgTypeCentr == utilities.Granted {
		fmt.Println("MESS REPLY !!!!!! ")

		myNode.Waiting = false
		myNode.ChanAcquireLock <- true

		//ora il processo puo entrare in CS
	}

	return nil
}

func checkAcks(process *utilities.NodeInfo) {

	//todo: quando azzero lista ReplyProSet ?????
	//date := time.Now().Format("15:04:05.000")
	fmt.Println("SONO IN checkAcks")

	fmt.Println("process.GetReplyProSet().Len() ==== ", process.GetReplyProSet().Len())
	fmt.Println("peers.Len()-1 ==== ", peers.Len()-1)
	fmt.Println("len(scalarMap) ==== ", len(myNode.ScalarMap))

	if process.GetReplyProSet().Len() == peers.Len()-1 && len(myNode.ScalarMap) > 0 {
		fmt.Println("HO RICEVUTO 2 MSG REPLY")

		//prendo il primo mess nella mappa per vedere se è il mio, ossia guardo ID sender
		msg := utilities.GetFirstElementMap(myNode.ScalarMap)
		fmt.Println("MSG IN CHECK ACKS ===", msg)

		if msg.Sender == myUsername {
			//il primo msg in lista e' il mio, quindi posso accedere in CS
			process.Waiting = false
			process.ChanAcquireLock <- true
		}
		//msg := heap.Pop(dl.msgHeap).(msgp3.Message)
		//vado qui quando il nodo ha ricevuto tutti i msg reply
		//dl.logger.Printf(date+" =====[1]=======,pop Front message(%v) to see whether is itselft.", msg.String())
		/*
			if msg.Sender == dl.nodeID && msg.TS == dl.requestTS {
				dl.logger.Printf(date+" lock(%v) has been notified.", dl.nodeID)
				dl.waiting = false
				heap.Push(dl.msgHeap, msg) // MUST push it back.
				dl.chanAcquireLock <- true
			} else {
				heap.Push(dl.msgHeap, msg) // MUST push it back.
			}

		*/
	}
}
