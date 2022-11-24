package ricartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var numMsg int
var logPaths *list.List
var numSender int

func ExecuteTestPeer(peer *RApeer, num int) {
	numSender = num
	logPaths = list.New().Init()

	fmt.Println("sto in ExecuteTestPeer")
	MyRApeer = *peer
	peerCnt = MyRApeer.PeerList.Len()

	if numSender == 1 && MyRApeer.ID == 0 {
		fmt.Println("mando il msg")
		SendRicart(&MyRApeer)
	}
	if numSender == 2 && (MyRApeer.ID == 0 || MyRApeer.ID == 1) {
		fmt.Println("mando il msg")
		SendRicart(&MyRApeer)
	} else {
		fmt.Println("sleep")
		time.Sleep(time.Minute)
	}

	fmt.Println(" ####################### TEST #############################")

	for i := 0; i < num; i++ {
		//fmt.Println(i)
		LogPath := "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(i) + ".log"
		logPaths.PushBack(LogPath)
	}

	if MyRApeer.ID == 1 {
		testNoStarvation()
	}

}

func testNoStarvation() {
	var csEntries int

	for e := logPaths.Front(); e != nil; e = e.Next() {

		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {
			//line := fileScanner.Text()

			fmt.Println(fileScanner.Text())
			if strings.Contains(fileScanner.Text(), "enters the critical section") {
				//fmt.Println("CONTIENE !!!!! ")
				csEntries++
			}
		}
		if numSender == 1 {
			break
		}
		//fmt.Println("\n---------------------------------\n\n")
	}
	//fmt.Println("csEntries == ", csEntries)

	if csEntries == numSender {
		fmt.Println(" === TEST NO STARVATION: PASSED !!")
	} else {
		fmt.Println(" === TEST NO STARVATION: FAILED !!")

	}
}
