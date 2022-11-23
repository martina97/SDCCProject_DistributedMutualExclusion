package tokenAsking

import (
	"fmt"
)

var num_msg int

func ExecuteTestPeer(peer *TokenPeer, numSender int) {
	fmt.Println("sto in ExecuteTestPeer")
	myPeer = *peer

	if numSender == 1 && myPeer.ID == 1 {
		fmt.Println("mando il msg")
	} else {
		fmt.Println("sleep")
	}

}

func ExecuteTestCoordinator(coordinator *Coordinator, numSender int) {
	fmt.Println("sto in ExecuteTestCoordinator")

	myCoordinator = *coordinator

	//aspetta finche il numero di token msg ricevuti è pari a numSender
	//Wait connection
	for num_msg < numSender { //todo: mettere 3 , anche sotto
		ch := <-Connection
		if ch == true {
			num_msg++
		}
	}
	fmt.Println("sto qua")
	Wg.Add(-numSender)
	fmt.Println("sto qua2")
}
