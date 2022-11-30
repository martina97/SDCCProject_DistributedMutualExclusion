package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	peers         *list.List
	myCoordinator Coordinator
	verbose       bool
)

func main() {

	flag.BoolVar(&verbose, "v", utilities.Verbose, "use this flag to get verbose info on messages")
	flag.Parse()

	if verbose {
		fmt.Println("VERBOSE FLAG ON")
	}

	peers = list.New()

	utilities.Registration(peers, utilities.ServerPort, "coordinator")
	fmt.Println("Registration completed successfully!")
	myCoordinator = *NewCoordinator()

	for e := peers.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.ServerPort))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go HandleConnectionCoordinator(connection)

	}
}

func HandleConnectionCoordinator(conn net.Conn) error {
	defer conn.Close()

	msg := new(tokenAsking.Message)
	dec := gob.NewDecoder(conn)
	err := dec.Decode(msg)
	utilities.CheckError(err, "error decoding message")

	/*
		if utilities.Test {
			go checkNumberToken()
		}

	*/

	if msg.MsgType == tokenAsking.Request {
		myCoordinator.mutex.Lock()
		err := tokenAsking.WriteMsgToFile("receive", *msg, myCoordinator.LogPath, true)
		utilities.CheckError(err, "error writing message")

		//time.Sleep(time.Second * 15)

		//devo controllare se è eleggibile!
		if tokenAsking.IsEligible(myCoordinator.VC, msg.VC, msg.Sender) && myCoordinator.HasToken {
			//invio token al processo e aggiorno il VC[i] del coordinatore, ossia incremento di 1 il valore relativo al processo
			myCoordinator.VC[msg.Sender]++
			sendToken(msg.Sender)
			myCoordinator.HasToken = false
			utilities.WriteVCInfoToFile(myCoordinator.LogPath, "coordinator", tokenAsking.ToString(myCoordinator.VC))
			utilities.WriteInfosToFile("gives token to "+msg.Sender, myCoordinator.LogPath, "coordinator")
		} else {
			//metto il msg in coda
			myCoordinator.ReqList.PushBack(msg)
		}
		myCoordinator.mutex.Unlock()
	}
	if msg.MsgType == tokenAsking.Token {
		myCoordinator.mutex.Lock()

		myCoordinator.numTokenMsgs++
		err := tokenAsking.WriteMsgToFile("receive", *msg, myCoordinator.LogPath, true)
		utilities.CheckError(err, "error writing message")

		myCoordinator.HasToken = true

		if myCoordinator.ReqList.Front() != nil {
			e := myCoordinator.ReqList.Front() //primo msg in coda
			pendingMsg := myCoordinator.ReqList.Front().Value.(*tokenAsking.Message)

			//vedo se il msg è eleggibile, e se sì invio msg con il token al sender del pendingMsg
			if tokenAsking.IsEligible(myCoordinator.VC, pendingMsg.VC, pendingMsg.Sender) {
				sendToken(pendingMsg.Sender)
				myCoordinator.HasToken = false
				utilities.WriteInfosToFile("gives token to "+pendingMsg.Sender, myCoordinator.LogPath, "coordinator")
				myCoordinator.ReqList.Remove(e)

				myCoordinator.VC[pendingMsg.Sender]++
				utilities.WriteVCInfoToFile(myCoordinator.LogPath, "coordinator", tokenAsking.ToString(myCoordinator.VC))
			}
		}
		/*
			if utilities.Test {
				Connection <- true
				Wg.Add(1)
			}

		*/

		myCoordinator.mutex.Unlock()
	}
	return nil
}

func sendToken(receiver string) {

	for e := peers.Front(); e != nil; e = e.Next() {
		dest := e.Value.(utilities.NodeInfo)
		if dest.Username == receiver {

			date := time.Now().Format(utilities.DateFormat)
			msg := tokenAsking.NewTokenMessage(date, "coordinator", receiver, myCoordinator.VC)
			utilities.SleepRandInt()

			peerConn := dest.Address + ":" + dest.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			utilities.CheckError(err, "Send response error on Dial")

			enc := gob.NewEncoder(conn)
			enc.Encode(msg)
			err = tokenAsking.WriteMsgToFile("send", *msg, myCoordinator.LogPath, true)
			utilities.CheckError(err, "error writing msg")
		}
	}

}
