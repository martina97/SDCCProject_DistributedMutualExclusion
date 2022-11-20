package tokenAsking

import "fmt"

var myPeer TokenPeer

func SendRequest(peer *TokenPeer) {
	if myPeer.Username == "" { //vuol dire che non ho ancora inizializzato il peer
		fmt.Println("sto in SendRequest --- RA_PEER VUOTA")
		myPeer = *peer

	} else {
		fmt.Println("sto in SendRequest --- RA_PEER NON VUOTA")
	}
	fmt.Println("myPeer.Coordinator= ", myPeer.Coordinator)

}
