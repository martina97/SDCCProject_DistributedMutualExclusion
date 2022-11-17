package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/RicartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

//send msg (o request o ack o release)
func sendMessage() error {
	switch algorithm {
	case "Centralized":
		sendCentralized()
	case "Lamport":
		sendLamport()
	case "RicartAgrawala":
		RicartAgrawala.SendRicart(myRApeer)
	}
	return nil
}

func sendCentralized() error {

	//myNode.SetReplyProSet(list.New())	//in questa lista metto il msg reply
	mu := myNode.GetMutex()
	mu.Lock()

	var res utilities.Result_file

	fmt.Println("SONO IN sendCentralized --- PORTA == ", myNode.Port)

	//only peer are destination of msgs

	//open connection whit peer
	//peerConn := dest.Address + ":" + dest.Port
	// todo: invece che peerConn mettere addr
	peerConn := utilities.Server_addr + ":" + strconv.Itoa(utilities.Server_port)

	conn, err := rpc.Dial("tcp", peerConn)
	defer conn.Close()
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	date := time.Now().Format(utilities.DATE_FORMAT)

	//msg := *utilities.NewEnterMsg(myNode, date)
	msg := *utilities.NewEnterMsg(myNode, date)
	fmt.Println("msg ==== ", msg)

	//call procedure
	log.Printf("Call to coordinator peer")
	err = conn.Call("Utility.CentralizedSincro", &msg, &res)
	if err != nil {
		log.Fatal("Error save_registration procedure: ", err)
	}

	myNode.Waiting = true

	mu.Unlock()

	<-myNode.ChanAcquireLock //il processo sta in attesa finche non riceve reply!!
	fmt.Println("dopo  <-myNode.ChanAcquireLock")

	return nil
}

func sendLamport() {

	myNode.SetReplyProSet(list.New())
	// tale lista serve a mettere i msg di reply per poi controllare che sono arrivati tutti
	// TODO: invece che lista basta semplicemente un contatore?!?!

	//lock := myNode.LockInfo
	mu := myNode.GetMutex()
	mu.Lock()
	//for range msgs {
	//increment local clock
	//incrementClock(&scalarClock, myID)

	utilities.IncrementTS(&myNode.TimeStamp)
	fmt.Println("------------------------------------------------------------- timestamp  ==", myNode.TimeStamp)

	/*
		seqNum := []uint64{}
		seqNum = append(seqNum, getValueClock(&scalarClock)[0])

	*/
	//date := time.Now().Format("2006/01/02 15:04:05")
	date := time.Now().Format(utilities.DATE_FORMAT)

	msg := *utilities.NewRequest2(myUsername, date, myNode.TimeStamp)

	fmt.Println("IL MESSAGGIO E' ====", msg)
	//fmt.Println("ID MESSAGGIO E' ====", msg.MsgID)
	fmt.Println("MsgType MESSAGGIO E' ====", msg.MsgType)
	fmt.Println("Sender MESSAGGIO E' ====", msg.Sender)
	fmt.Println("Date MESSAGGIO E' ====", msg.Date)
	fmt.Println("timeStamp MESSAGGIO E' ====", msg.TS)
	sendRequest(msg)

	myNode.Waiting = true

	mu.Unlock()

	utilities.WriteInfoToFile(myID, " wait all peer reply messages.", false)
	/*
		f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		//save new address on file
		date = time.Now().Format("15:04:05.000")


		_, err = f.WriteString(date + " lock(" + strconv.Itoa(myID) + ") wait all peer reply messages.\n")

		_, err = f.WriteString("\n")
		err = f.Sync()

	*/

	<-myNode.ChanAcquireLock

	utilities.WriteInfoToFile(myID, " receive all peer reply messages successfully.", false)
	/*
		date = time.Now().Format("15:04:05.000")

		_, err = f.WriteString(date + " lock(" + strconv.Itoa(myID) + ")receive all peer reply messages successfully.\n") //todo:invece che lock scrivere processo

		_, err = f.WriteString("\n")
		err = f.Sync()

	*/

	//ho ricevuto tutti msg reply, ora entro in cs
	fmt.Println("lista di msg in coda ==", myNode.ScalarMap)
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
	utilities.IncrementTS(&myNode.TimeStamp)

	date := time.Now().Format(utilities.DATE_FORMAT)

	releaseMsg := *utilities.NewRelease(myUsername, date, myNode.TimeStamp)
	utilities.WriteTSInfoToFile(myID, myNode.TimeStamp, algorithm)

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
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

			releaseMsg.Receiver = dest.Username
			/*
				f, err := os.OpenFile("/docker/node_volume/process_"+strconv.Itoa(myID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
				if err != nil {
					log.Fatalf("error opening file: %v", err)
				}
				//save new address on file
				date = time.Now().Format("15:04:05.000")
				_, err = f.WriteString("[" + date + "] : Send " + releaseMsg.MessageToString() + " to peer " + strconv.Itoa(dest.ID))
				_, err = f.WriteString("\n")
				err = f.Sync()
				if err != nil {
					return err
				}

			*/

			err = utilities.WriteMsgToFile(&myNode, "Send", releaseMsg, dest.ID, myNode.TimeStamp)

			if err != nil {
				return err
			}
		}
	}

	//elimino primo msg da lista
	utilities.RemoveFirstElementMap(myNode.ScalarMap)
	fmt.Println("ora la mappa ===", myNode.ScalarMap)
	return nil
}

func sendRequest(msg utilities.Message) error {

	utilities.WriteTSInfoToFile(myID, myNode.TimeStamp, algorithm)
	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
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

			msg.Receiver = dest.Username

			err = utilities.WriteMsgToFile(&myNode, "Send", msg, dest.ID, myNode.TimeStamp)
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
				_, err = f.WriteString("[" + date + "] : Send " + msg.MessageToString() + " to peer " + strconv.Itoa(dest.ID))
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

	utilities.AppendHashMap2(myNode.ScalarMap, msg)
	fmt.Println(" ------------------------------------------ STO QUA 3 ----------------------------")

	/*
		for e := lista(); e != nil; e = e.Next() {
			item := e.Value.(Message)

		}
		fmt.Println("LISTA DEL PEER SENDER ==", lista)

	*/

	fmt.Println("MAPPA SENDER ====", myNode.ScalarMap)

	return nil
}

func sendAck(msg *utilities.Message) error {
	// mando ack al peer con id msg.receiver
	utilities.WriteTSInfoToFile(myID, myNode.TimeStamp, algorithm)

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		if dest.Username == msg.Receiver {
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
			date := time.Now().Format(utilities.DATE_FORMAT)
			_, err = f.WriteString("[" + date + "] : Send" + msg.MessageToString("send") + " to peer " + strconv.Itoa(dest.ID))
			_, err = f.WriteString("\n")
			err = f.Sync()
			if err != nil {
				return err
			}
		}
	}
	return nil

}
