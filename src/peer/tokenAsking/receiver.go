package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"net"
)

//Save message
func HandleConnectionCoordinator(conn net.Conn, coordinator *Coordinator) error {
	defer conn.Close()

	fmt.Println("sto in HandleConnectionCoordinator dentro tokenAsking package")
	fmt.Println("coordinator === ", coordinator)

	if myCoordinator.Username == "" { //vuol dire che non ho ancora inizializzato il coordinatore
		fmt.Println("sto in HandleConnectionCoordinator --- coordinator VUOTA")
		myCoordinator = *coordinator

	} else {
		fmt.Println("sto in HandleConnectionCoordinator --- coordinator NON VUOTA")
	}
	fmt.Println("myCoordinator = ", myCoordinator)

	//myCoordinator = *coordinator

	msg := new(Message)
	dec := gob.NewDecoder(conn)
	dec.Decode(msg)
	fmt.Println("il msg == ", msg.ToString("receive"))

	if msg.MsgType == Request {
		fmt.Println("sto dentro if")
		myCoordinator.mutex.Lock()
		WriteMsgToFile("receive", *msg, true)
		fmt.Println("ricevo ", msg, "e hastoken == ", myCoordinator.HasToken)
		//devo controllare se è eleggibile!
		if utilities.IsEligible(myCoordinator.VC, msg.VC, msg.Sender) && myCoordinator.HasToken {
			fmt.Println("msg eleggibile!")
			//update VC
			myCoordinator.VC[msg.Sender]++
			fmt.Println("vc coord = ", myCoordinator.VC)

			//invio token al processo e aggiorno il VC[i] del coordinatore, ossia incremento di 1 il valore relativo al processo
			sendToken(msg.Sender, true)
			myCoordinator.HasToken = false
			WriteVCInfoToFile(true)
			fmt.Println("il coordinatore non ha piu token! ")
			WriteInfosToFile("gives token to "+msg.Sender, true)
			fmt.Println("hasToken ==", myCoordinator.HasToken)

		} else {
			fmt.Println("msg non eleggibile")
			//metto il msg in coda
			myCoordinator.ReqList.PushBack(msg)
		}
		myCoordinator.mutex.Unlock()
	}
	if msg.MsgType == Token {
		fmt.Println("msg Type === TOKEN, msg = ", msg)
		myCoordinator.mutex.Lock()
		WriteMsgToFile("receive", *msg, true)
		myCoordinator.HasToken = true
		//WriteInfosToFile("gets the token.")
		fmt.Println("STO QUA!")

		if myCoordinator.ReqList.Front() != nil {
			fmt.Println("in coda c'è :", myCoordinator.ReqList.Front().Value)
			e := myCoordinator.ReqList.Front()
			pendingMsg := myCoordinator.ReqList.Front().Value.(*Message)
			fmt.Println("pendingMsg :", pendingMsg)
			fmt.Println("in coda c'è :", myCoordinator.ReqList.Front().Value)

			//vedo se il msg è eleggibile, e se sì invio msg con il token al sender del pendingMsg
			if utilities.IsEligible(myCoordinator.VC, pendingMsg.VC, pendingMsg.Sender) {
				fmt.Println("posso inviare il token al peer")
				sendToken(pendingMsg.Sender, true)
				myCoordinator.HasToken = false
				WriteInfosToFile("gives token to "+pendingMsg.Sender, true)
				myCoordinator.ReqList.Remove(e)
				myCoordinator.VC[pendingMsg.Sender]++

				WriteVCInfoToFile(true)

				fmt.Println("req List == ", myCoordinator.ReqList)
			}

		} else {
			fmt.Println("coda richieste pendenti vuota!")
		}

		myCoordinator.mutex.Unlock()

	}

	return nil
}

func HandleConnectionPeer(conn net.Conn, peer *TokenPeer) error {

	fmt.Println("sto in HandleConnection dentro tokenAsking package")
	fmt.Println("peer === ", peer)

	if myPeer.Username == "" {
		fmt.Println("peer VUOTA")
		myPeer = *peer
		//peerCnt = MyRApeer.PeerList.Len()
	} else {
		fmt.Println("peer NON VUOTA")
	}

	fmt.Println("peer == ", myPeer)
	defer conn.Close()

	msg := new(Message)
	dec := gob.NewDecoder(conn)
	dec.Decode(msg)
	fmt.Println("il msg == ", msg.ToString("receive"))
	if msg.MsgType == ProgramMessage {
		myPeer.mutex.Lock()
		//update VC !
		utilities.UpdateVC(myPeer.VC, msg.VC)
		WriteMsgToFile("receive", *msg, false)

		myPeer.mutex.Unlock()
	}
	if msg.MsgType == Token {
		// ho il token !
		myPeer.mutex.Lock()
		WriteMsgToFile("receive", *msg, false)
		myPeer.HasToken <- true
		myPeer.mutex.Unlock()

	}

	return nil
}
