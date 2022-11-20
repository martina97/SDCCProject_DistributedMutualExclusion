package tokenAsking

import "fmt"

var myPeer TokenPeer
var myCoordinator Coordinator

func SendRequest(peer *TokenPeer, coordinator *Coordinator) {
	if myPeer == (TokenPeer{}) {
		fmt.Println("sto in SendRequest --- RA_PEER VUOTA")
		myPeer = *peer

	} else {
		fmt.Println("sto in SendRequest --- RA_PEER NON VUOTA")
	}
	myCoordinator = *coordinator
	fmt.Println("myCoordinator = ", myCoordinator)

}
