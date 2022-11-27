package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"net"
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
		err := WriteMsgToFile("receive", *msg, true)
		utilities.CheckError(err, "error writing message")

		//time.Sleep(time.Second * 15)

		//devo controllare se è eleggibile!
		if IsEligible(myCoordinator.VC, msg.VC, msg.Sender) && myCoordinator.HasToken {
			//invio token al processo e aggiorno il VC[i] del coordinatore, ossia incremento di 1 il valore relativo al processo
			myCoordinator.VC[msg.Sender]++
			sendToken(msg.Sender, true)
			myCoordinator.HasToken = false
			WriteVCInfoToFile(true)
			WriteInfosToFile("gives token to "+msg.Sender, true)
		} else {
			//metto il msg in coda
			myCoordinator.ReqList.PushBack(msg)
		}
		myCoordinator.mutex.Unlock()
	}
	if msg.MsgType == Token {
		myCoordinator.mutex.Lock()

		myCoordinator.numTokenMsgs++
		err := WriteMsgToFile("receive", *msg, true)
		utilities.CheckError(err, "error writing message")

		myCoordinator.HasToken = true

		if myCoordinator.ReqList.Front() != nil {
			e := myCoordinator.ReqList.Front() //primo msg in coda
			pendingMsg := myCoordinator.ReqList.Front().Value.(*Message)

			//vedo se il msg è eleggibile, e se sì invio msg con il token al sender del pendingMsg
			if IsEligible(myCoordinator.VC, pendingMsg.VC, pendingMsg.Sender) {
				sendToken(pendingMsg.Sender, true)
				myCoordinator.HasToken = false
				WriteInfosToFile("gives token to "+pendingMsg.Sender, true)
				myCoordinator.ReqList.Remove(e)

				myCoordinator.VC[pendingMsg.Sender]++
				WriteVCInfoToFile(true)
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
		err := WriteMsgToFile("receive", *msg, false)
		utilities.CheckError(err, "error writing msg")
		myPeer.mutex.Unlock()
	}

	if msg.MsgType == Token {
		// ho il token !
		myPeer.mutex.Lock()
		err := WriteMsgToFile("receive", *msg, false)
		utilities.CheckError(err, "error writing msg")
		myPeer.HasToken <- true
		myPeer.mutex.Unlock()
	}

	return nil
}
