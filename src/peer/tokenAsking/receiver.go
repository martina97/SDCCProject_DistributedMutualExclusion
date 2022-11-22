package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"net"
)

//Save message
func HandleConnectionCoordinator(conn net.Conn, coordinator *Coordinator) error {

	fmt.Println("sto in HandleConnectionCoordinator dentro tokenAsking package")
	fmt.Println("coordinator === ", coordinator)
	myCoordinator = *coordinator

	defer conn.Close()
	msg := new(Message)
	dec := gob.NewDecoder(conn)
	dec.Decode(msg)
	fmt.Println("il msg == ", msg.ToString("receive"))

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
		WriteMsgToFile("receive", *msg)

		myPeer.mutex.Unlock()
	}

	return nil
}
