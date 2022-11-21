package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

var myPeer TokenPeer
var myCoordinator Coordinator

func SendRequest(peer *TokenPeer) {
	if myPeer.Username == "" { //vuol dire che non ho ancora inizializzato il peer
		fmt.Println("sto in SendRequest --- RA_PEER VUOTA")
		myPeer = *peer

	} else {
		fmt.Println("sto in SendRequest --- RA_PEER NON VUOTA")
	}
	fmt.Println("myPeer.Coordinator= ", myPeer.Coordinator)

	myPeer.mutex.Lock()
	//incremento Vector Clock!!!
	fmt.Println("myTokenPeer.VC =", myPeer.VC)
	fmt.Println("incremento VC")
	utilities.IncrementVC(myPeer.VC, myPeer.Username)

	fmt.Println("myTokenPeer.VC =", myPeer.VC)

}
