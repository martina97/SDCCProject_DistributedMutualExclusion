package tokenAsking

import "fmt"

func ExecuteTestPeer(peer *TokenPeer) {
	if myPeer.Username == "" { //vuol dire che non ho ancora inizializzato il peer
		fmt.Println("sto in ExecuteTestPeer --- RA_PEER VUOTA")
		myPeer = *peer

	} else {
		fmt.Println("sto in ExecuteTestPeer --- RA_PEER NON VUOTA")
	}

}

func ExecuteTestCoordinator(coordinator *Coordinator) {
	if myCoordinator.Username == "" { //vuol dire che non ho ancora inizializzato il coordinatore
		fmt.Println("sto in ExecuteTestCoordinator --- coordinator VUOTA")
		myCoordinator = *coordinator

	} else {
		fmt.Println("sto in ExecuteTestCoordinator --- coordinator NON VUOTA")
	}
}