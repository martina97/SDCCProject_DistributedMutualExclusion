package main

import (
	"SDCCProject_DistributedMutualExclusion/app/utilities"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

//send msg (o request o ack o release)
func sendMessages() error {
	sendLamport()
	return nil
}

func sendLamport() {

	myProcess.SetReplyProSet(list.New())
	// tale lista serve a mettere i msg di reply per poi controllare che sono arrivati tutti
	// TODO: invece che lista basta semplicemente un contatore?!?!

	//lock := myProcess.LockInfo
	mu := myProcess.GetMutex()
	mu.Lock()
	//for range msgs {
	//increment local clock
	//incrementClock(&scalarClock, myID)

	utilities.IncrementTS(&timeStamp)
	fmt.Println("------------------------------------------------------------- timestamp  ==", timeStamp)

	/*
		seqNum := []uint64{}
		seqNum = append(seqNum, getValueClock(&scalarClock)[0])

	*/
	//date := time.Now().Format("2006/01/02 15:04:05")
	date := time.Now().Format("15:04:05.000")

	msg := *utilities.NewRequest2(myID, date, timeStamp)

	fmt.Println("IL MESSAGGIO E' ====", msg)
	//fmt.Println("ID MESSAGGIO E' ====", msg.MsgID)
	fmt.Println("MsgType MESSAGGIO E' ====", msg.MsgType)
	fmt.Println("Sender MESSAGGIO E' ====", msg.Sender)
	fmt.Println("Date MESSAGGIO E' ====", msg.Date)
	fmt.Println("timeStamp MESSAGGIO E' ====", msg.TS)
	sendRequest(msg)

	myProcess.Waiting = true

	mu.Unlock()

	utilities.WriteInfoToFile(myID, " wait all node reply messages.", false)
	/*
		f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		//save new address on file
		date = time.Now().Format("15:04:05.000")


		_, err = f.WriteString(date + " lock(" + strconv.Itoa(myID) + ") wait all node reply messages.\n")

		_, err = f.WriteString("\n")
		err = f.Sync()

	*/

	<-myProcess.ChanAcquireLock

	utilities.WriteInfoToFile(myID, " receive all node reply messages successfully.", false)
	/*
		date = time.Now().Format("15:04:05.000")

		_, err = f.WriteString(date + " lock(" + strconv.Itoa(myID) + ")receive all node reply messages successfully.\n") //todo:invece che lock scrivere processo

		_, err = f.WriteString("\n")
		err = f.Sync()

	*/

	//ho ricevuto tutti msg reply, ora entro in cs
	fmt.Println("lista di msg in coda ==", scalarMap)
	fmt.Println("entro in CS")
	utilities.WriteInfoToFile(myID, " entered the critical section at ", true)
	time.Sleep(time.Minute / 2) //todo: invece che sleep mettere file condiviso
	utilities.WriteInfoToFile(myID, " exited the critical section at ", true)

	//log.Writer()

	/*
			date = time.Now().Format("15:04:05.000")
			_, err = f.WriteString("process " + strconv.Itoa(myID) + " entered the critical section at " + date)
			_, err = f.WriteString("\n")
			err = f.Sync()



		date = time.Now().Format("15:04:05.000")
		_, err = f.WriteString("process " + strconv.Itoa(myID) + " exited the critical section at " + date)
		_, err = f.WriteString("\n")
		err = f.Sync()
		fmt.Println("uscito da CS")

	*/

	//lascio CS e mando msg release a tutti
	sendRelease()

	//prepare msg to send
	//var msg utilities.Message
	//msg.MsgID = "prova"
	/*
		msg.MsgType = utilities.Request
		msg.SeqNum = append(msg.SeqNum, getValueClock(&scalarClock)[0])
		msg.Date = time.Now().Format("2006/01/02 15:04:05")
		//msg.Text = text
		msg.Sender = myID

	*/
}

func sendRelease() error {
	//incremento timestamp
	utilities.IncrementTS(&timeStamp)

	date := time.Now().Format("15:04:05.000")

	releaseMsg := *utilities.NewRelease(myID, date, timeStamp)
	utilities.WriteTSInfoToFile(myID, timeStamp)

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.Process)
		//only peer are destination of msgs
		if dest.Type == utilities.Peer && dest.ID != myID { //non voglio mandarlo a me stesso

			//open connection whit peer
			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(releaseMsg)

			releaseMsg.Receiver = dest.ID
			/*
				f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
				if err != nil {
					log.Fatalf("error opening file: %v", err)
				}
				//save new address on file
				date = time.Now().Format("15:04:05.000")
				_, err = f.WriteString("[" + date + "] : Send " + releaseMsg.MessageToString() + " to node " + strconv.Itoa(dest.ID))
				_, err = f.WriteString("\n")
				err = f.Sync()
				if err != nil {
					return err
				}

			*/

			err = utilities.WriteMsgToFile(&myProcess, "Send", releaseMsg, dest.ID, timeStamp)

			if err != nil {
				return err
			}
		}
	}

	//elimino primo msg da lista
	utilities.RemoveFirstElementMap(scalarMap)
	fmt.Println("ora la mappa ===", scalarMap)
	return nil
}

func sendRequest(msg utilities.Message) error {

	utilities.WriteTSInfoToFile(myID, timeStamp)
	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.Process)
		//only peer are destination of msgs
		if dest.Type == utilities.Peer && dest.ID != myID { //non voglio mandarlo a me stesso

			//open connection whit peer
			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			msg.Receiver = dest.ID

			err = utilities.WriteMsgToFile(&myProcess, "Send", msg, dest.ID, timeStamp)
			if err != nil {
				return err
			}

			/*
				f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
				if err != nil {
					log.Fatalf("error opening file: %v", err)
				}
				//save new address on file
				date := time.Now().Format("15:04:05.000")
				_, err = f.WriteString("[" + date + "] : Send " + msg.MessageToString() + " to node " + strconv.Itoa(dest.ID))
				_, err = f.WriteString("\n")
				err = f.Sync()
				if err != nil {
					return err
				}

			*/
		}
	}
	//una volta inviato il msg, lo salvo nella coda locale del peer sender
	fmt.Println(" ------------------------------------------ STO QUA 2 ----------------------------")

	utilities.AppendHashMap2(scalarMap, msg)
	fmt.Println(" ------------------------------------------ STO QUA 3 ----------------------------")

	/*
		for e := lista(); e != nil; e = e.Next() {
			item := e.Value.(Message)

		}
		fmt.Println("LISTA DEL PEER SENDER ==", lista)

	*/

	fmt.Println("MAPPA SENDER ====", scalarMap)

	return nil
}

func sendAck(msg *utilities.Message) error {
	// mando ack al peer con id msg.receiver
	utilities.WriteTSInfoToFile(myID, timeStamp)

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.Process)
		if dest.ID == msg.Receiver {
			//open connection whit peer
			peer_conn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peer_conn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			//save new address on file
			date := time.Now().Format("15:04:05.000")
			_, err = f.WriteString("[" + date + "] : Send" + msg.MessageToString("send") + " to node " + strconv.Itoa(dest.ID))
			_, err = f.WriteString("\n")
			err = f.Sync()
			if err != nil {
				return err
			}
		}
	}
	return nil

}
