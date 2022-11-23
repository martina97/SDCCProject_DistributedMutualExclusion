package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var num_msg int
var logPaths *list.List

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

	logPaths = list.New().Init()
	//logPaths.Init()
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

	for i := 0; i < utilities.MAXPEERS; i++ {
		if i != myCoordinator.ID {
			//fmt.Println(i)
			LogPath := "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(i) + ".log"
			logPaths.PushBack(LogPath)
			//fmt.Println(LogPath)
		}
	}
	checkSafety()

	/*
		ora posso controllare i vari file di log!!
		1 coordinator.log
		n-1 peer_n.log
	*/

	fileScanner := getFileSplit(myCoordinator.LogPath)
	for fileScanner.Scan() {
		//line := fileScanner.Text()

		fmt.Println(fileScanner.Text())
	}

	//f.Close()
}

func checkSafety() {

	for e := logPaths.Front(); e != nil; e = e.Next() {

	}

}

func getFileSplit(path string) *bufio.Scanner {
	//provo a farlo con coordinator.log
	f, err := os.OpenFile(myCoordinator.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	fmt.Println("sto qua3")

	fileScanner := bufio.NewScanner(f)
	fmt.Println("sto qua4")

	fileScanner.Split(bufio.ScanLines)
	fmt.Println("sto qua5")
	return fileScanner
}
