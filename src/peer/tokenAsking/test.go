package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"bufio"
	"fmt"
	"time"
)

var num_msg int

func ExecuteTestPeer(peer *TokenPeer, numSender int) {
	fmt.Println("sto in ExecuteTestPeer")
	myPeer = *peer

	if numSender == 1 && myPeer.ID == 1 {
		fmt.Println("mando il msg")
		SendRequest(&myPeer)
	} else {
		fmt.Println("sleep")
		time.Sleep(time.Minute / 2)
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

	checkSafety()

	/*
		ora posso controllare i vari file di log!!
		1 coordinator.log
		n-1 peer_n.log
	*/

	//provo a farlo con coordinator.log
	f := openFile(true)
	fmt.Println("sto qua3")

	fileScanner := bufio.NewScanner(f)
	fmt.Println("sto qua4")

	fileScanner.Split(bufio.ScanLines)
	fmt.Println("sto qua5")

	for fileScanner.Scan() {
		//line := fileScanner.Text()

		fmt.Println(fileScanner.Text())
	}

	f.Close()
}

func checkSafety() {

	//devo aprire i file dei processi
	for e := myCoordinator.PeerList.Front(); e != nil; e = e.Next() {
		peer := e.Value.(*utilities.NodeInfo)
		if peer.Username != utilities.COORDINATOR {
			fmt.Println(peer.LogPath)
		}
	}
}
