package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func HandleConnectionCoordinator(conn net.Conn, coordinator *Coordinator) error {
	defer conn.Close()

	if myCoordinator.Username == "" { //vuol dire che non ho ancora inizializzato il coordinatore
		myCoordinator = *coordinator
	}

	msg := new(Message)
	dec := gob.NewDecoder(conn)
	err := dec.Decode(msg)
	utilities.CheckError(err, "error decoding message")

	if msg.MsgType == Request {
		myCoordinator.mutex.Lock()
		err := WriteMsgToFile("receive", *msg, myCoordinator.LogPath, true)
		utilities.CheckError(err, "error writing message")

		//time.Sleep(time.Second * 15)

		//devo controllare se è eleggibile!
		if IsEligible(myCoordinator.VC, msg.VC, msg.Sender) && myCoordinator.HasToken {
			//invio token al processo e aggiorno il VC[i] del coordinatore, ossia incremento di 1 il valore relativo al processo
			myCoordinator.VC[msg.Sender]++
			sendToken(msg.Sender, true)
			myCoordinator.HasToken = false
			utilities.WriteVCInfoToFile(myCoordinator.LogPath, "coordinator", ToString(myCoordinator.VC))
			utilities.WriteInfosToFile("gives token to "+msg.Sender, myCoordinator.LogPath, "coordinator")
		} else {
			//metto il msg in coda
			myCoordinator.ReqList.PushBack(msg)
		}
		myCoordinator.mutex.Unlock()
	}
	if msg.MsgType == Token {
		myCoordinator.mutex.Lock()

		myCoordinator.numTokenMsgs++
		err := WriteMsgToFile("receive", *msg, myCoordinator.LogPath, true)
		utilities.CheckError(err, "error writing message")

		myCoordinator.HasToken = true

		if myCoordinator.ReqList.Front() != nil {
			e := myCoordinator.ReqList.Front() //primo msg in coda
			pendingMsg := myCoordinator.ReqList.Front().Value.(*Message)

			//vedo se il msg è eleggibile, e se sì invio msg con il token al sender del pendingMsg
			if IsEligible(myCoordinator.VC, pendingMsg.VC, pendingMsg.Sender) {
				sendToken(pendingMsg.Sender, true)
				myCoordinator.HasToken = false
				utilities.WriteInfosToFile("gives token to "+pendingMsg.Sender, myCoordinator.LogPath, "coordinator")
				myCoordinator.ReqList.Remove(e)

				myCoordinator.VC[pendingMsg.Sender]++
				utilities.WriteVCInfoToFile(myCoordinator.LogPath, "coordinator", ToString(myCoordinator.VC))
			}
		}
		if utilities.Test {
			Connection <- true
			Wg.Add(1)
		}
		myCoordinator.mutex.Unlock()
	}
	return nil
}

func HandleConnectionPeer(conn net.Conn, peer *TokenPeer) error {

	if myPeer.Username == "" {
		myPeer = *peer
	}
	defer conn.Close()

	msg := new(Message)
	dec := gob.NewDecoder(conn)
	err := dec.Decode(msg)
	utilities.CheckError(err, "error decoding msg")

	if msg.MsgType == ProgramMessage {
		myPeer.mutex.Lock()
		//update VC !
		UpdateVC(myPeer.VC, msg.VC)
		err := WriteMsgToFile("receive", *msg, myPeer.LogPath, false)
		utilities.CheckError(err, "error writing msg")
		myPeer.mutex.Unlock()
	}

	if msg.MsgType == Token {
		// ho il token !
		myPeer.mutex.Lock()
		err := WriteMsgToFile("receive", *msg, myPeer.LogPath, false)
		utilities.CheckError(err, "error writing msg")
		myPeer.HasToken = true
		myPeer.mutex.Unlock()
	}

	return nil
}

func checkHasToken() {

	fmt.Println("sto in checkHasToken")
	for !(myPeer.HasToken) {
		fmt.Println("sto in checkHasToken dentro for")

		time.Sleep(time.Second * 5)
	}
	fmt.Println("sto in checkHasToken fuori for")

	date := time.Now().Format(utilities.DateFormat)
	utilities.WriteInfosToFile("enters the critical section at "+date+".", myPeer.LogPath, myPeer.Username)
	time.Sleep(time.Minute / 2)
	date = time.Now().Format(utilities.DateFormat)

	utilities.WriteInfosToFile("exits the critical section at "+date+".", myPeer.LogPath, myPeer.Username)
	utilities.WriteInfosToFile("releases the token.", myPeer.LogPath, myPeer.Username)

	myPeer.mutex.Lock()

	//invio msg con il token al coordinatore
	sendToken("coordinator", false)
	myPeer.mutex.Unlock()

}
