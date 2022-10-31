package node

import (
	"SDCCProject_DistributedMutualExclusion/app/utilities"
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
		go handleConnection(connection)
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
	//ogni peer ha il suo clock scalare, e' var globale come MyProcess e myID

	time.Sleep(time.Minute / 2) //PRIMA DI AUMENTARE TS METTO SLEEP COSI PROVO A INVIARE 2 REQ INSIEME E VEDO CHE SUCCEDE

	utilities.UpdateTS(&MyProcess.TimeStamp, &msg.TS)

	//mutex := lock.GetMutex()
	mutex := MyProcess.GetMutex()
	if msg.MsgType == utilities.Request {
		/*
			quando ricevo una richiesta da un processo devo decidere se mandare ACK al processo oppure se voglio entrare in CS
		*/
		fmt.Println("MESS REQUEST !!!!!! ")
		fmt.Println("TIMESTAMP QUANDO RICEVO REQUEST ===", MyProcess.TimeStamp) //ho gia aggiornato il TS!!
		//fmt.Println("------------------------------------------------------------- DOPO RICEVUTO REQUEST --- > timestamp  ==", timeStamp)

		mutex.Lock()
		utilities.WriteMsgToFile(&MyProcess, "Receive", *msg, 0, MyProcess.TimeStamp)
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
		utilities.AppendHashMap2(MyProcess.ScalarMap, *msg)

		//QUA DEVO DECIDERE SE MANDARE ACK O REQUEST (msg REPLY O REQUEST)

		//scelta := openMenuRequest()
		//fmt.Println("HO SCELTO", scelta)

		//mando msg reply
		//date := time.Now().Format("17:06:04")
		//prima di mandare reply aggiorno di nuovo il TS !!
		utilities.IncrementTS(&MyProcess.TimeStamp)

		fmt.Println("------------------------------------------------------------- DOPO INVIATO REPLY --- > timestamp  ==", MyProcess.TimeStamp)
		date := time.Now().Format("15:04:05.000")
		replyMsg := utilities.NewReply(tmp, MyProcess.ID, msg.Sender, date, MyProcess.TimeStamp)
		sendAck(replyMsg)
		mutex.Unlock()
	}

	if msg.MsgType == utilities.Reply {
		fmt.Println("------------------------------------------------------------- DOPO RICEVUTO REPLY --- > timestamp  ==", MyProcess.TimeStamp)
		mutex.Lock()
		fmt.Println("TIMESTAMP QUANDO RICEVO Reply ===", MyProcess.TimeStamp)

		utilities.WriteMsgToFile(&MyProcess, "Receive", *msg, 0, MyProcess.TimeStamp)
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
		MyProcess.GetReplyProSet().PushBack(msg)
		//check ack
		checkAcks(&MyProcess) //controllo se ho ricevuto 2 msg reply, se si posso entrare in CS prendendo 1 elem nella lista
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
		utilities.WriteMsgToFile(&MyProcess, "Receive", *msg, 0, MyProcess.TimeStamp)

		//utilities.WriteTSInfoToFile(myID, timeStamp)

		utilities.RemoveFirstElementMap(MyProcess.ScalarMap)
		fmt.Println("---------------------------------	DOPO AVER RICEVUTO RELEASE mappa ===", MyProcess.ScalarMap)
		checkAcks(&MyProcess)
		mutex.Unlock()

	}

	fmt.Println("msg ricevuti -----", MyProcess.ScalarMap)
	/*
		for key, element := range scalarMap {
			fmt.Println("Key:", key, "=>", "Element:", element)
		}

	*/
	//fmt.Println("PRIMO MESS ==", utilities.GetFirstElementMap(scalarMap))

	return nil
}

func checkAcks(process *utilities.Process) {

	//todo: quando azzero lista ReplyProSet ?????
	//date := time.Now().Format("15:04:05.000")
	fmt.Println("SONO IN checkAcks")

	fmt.Println("process.GetReplyProSet().Len() ==== ", process.GetReplyProSet().Len())
	fmt.Println("peers.Len()-1 ==== ", peers.Len()-1)
	fmt.Println("len(scalarMap) ==== ", len(MyProcess.ScalarMap))

	if process.GetReplyProSet().Len() == peers.Len()-1 && len(MyProcess.ScalarMap) > 0 {
		fmt.Println("HO RICEVUTO 2 MSG REPLY")

		//prendo il primo mess nella mappa per vedere se è il mio, ossia guardo ID sender
		msg := utilities.GetFirstElementMap(MyProcess.ScalarMap)
		fmt.Println("MSG IN CHECK ACKS ===", msg)

		if msg.Sender == myID {
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
