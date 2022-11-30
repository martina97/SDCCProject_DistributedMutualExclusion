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

	MyRApeer = *peer
	peerCnt = MyRApeer.PeerList.Len()

	if numSender == 1 && MyRApeer.ID == 0 {
		SendRicart(&MyRApeer, verbose)
		<-MyRApeer.ChanStartTest
	}
	if numSender == 2 && (MyRApeer.ID == 0 || MyRApeer.ID == 1) {
		SendRicart(&MyRApeer, verbose)
		<-MyRApeer.ChanStartTest

	} else {
		time.Sleep(time.Minute)
	}

	fmt.Println(" ####################### TEST #############################")

	for i := 0; i < num; i++ {
		//fmt.Println(i)
		LogPath := "/docker/node_volume/ricartAgrawala/peer_" + strconv.Itoa(i) + ".log"
		logPaths.PushBack(LogPath)
	}

	if MyRApeer.ID == 1 {
		testMessageNumber()
		testNoStarvation()
		if numSender == 2 {
			testSafety()
		}
	}

}

func testNoStarvation() {
	var csEntries int

	for e := logPaths.Front(); e != nil; e = e.Next() {

		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {
			if strings.Contains(fileScanner.Text(), "enters the critical section") {
				csEntries++
			}
		}
	}

	if csEntries == numSender {
		fmt.Println(" === TEST NO STARVATION: PASSED !!")
	} else {
		fmt.Println(" === TEST NO STARVATION: FAILED !!")

	}
}

//solo se numSender = 2
func testSafety() {

	stringEnter := "enters the critical section at "
	stringExit := "exits the critical section at "

	var enterP1 time.Time
	var enterP2 time.Time
	var exitP1 time.Time
	var exitP2 time.Time
	var result bool
	index := 0

	for e := logPaths.Front(); e != nil; e = e.Next() {
		var enterDate time.Time
		var exitDate time.Time

		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {

			if strings.Contains(fileScanner.Text(), stringEnter) {
				enterDate = utilities.ConvertStringToDate(fileScanner.Text(), stringEnter)
			}
			if strings.Contains(fileScanner.Text(), stringExit) {
				exitDate = utilities.ConvertStringToDate(fileScanner.Text(), stringExit)
			}
		}
		if index == 0 {
			//prima iterazione
			enterP1 = enterDate
			exitP1 = exitDate
		} else {
			enterP2 = enterDate
			exitP2 = exitDate
		}

		index++

	}

	if enterP1.Before(enterP2) {
		result = exitP1.Before(enterP2)
	} else {
		result = exitP2.Before(enterP1)
	}
	if result {
		fmt.Println(" === TEST SAFETY: PASSED !!")
	} else {
		fmt.Println(" === TEST SAFETY: FAILED !!")

	}
}

func testMessageNumber() {
	/*
		2(N-1) messaggi per accedere alla CS:
			• N-1 messaggi di richiesta
			• N-1 messaggi di reply
	*/

	index := 0
	for e := logPaths.Front(); e != nil; e = e.Next() {
		numMsg := 0
		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {

			if strings.Contains(fileScanner.Text(), "sends Request message") ||
				strings.Contains(fileScanner.Text(), "receives Reply message") {
				numMsg++
			}
		}

		if numMsg == 2*(utilities.MAXPEERS-1) {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : PASSED !!!\n", index)
		} else {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : FAILED !!!\n", index)
		}
		index++

	}
}
