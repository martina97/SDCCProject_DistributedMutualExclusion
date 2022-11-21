package tokenAsking

import (
	"fmt"
	"net"
)

//Save message
func HandleConnectionCoordinator(conn net.Conn, coordinator *Coordinator) error {

	fmt.Println("sto in HandleConnectionCoordinator dentro tokenAsking package")
	fmt.Println("coordinator === ", coordinator)

	return nil
}

func HandleConnectionPeer(conn net.Conn, peer *TokenPeer) error {

	fmt.Println("sto in HandleConnectionCoordinator dentro tokenAsking package")
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

	return nil
}
