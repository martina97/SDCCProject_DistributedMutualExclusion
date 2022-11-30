package tokenAsking

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

func ExecuteTestPeer(peer *TokenPeer, num int) {
	numSender = num
	myPeer = *peer

	if numSender == 1 && myPeer.ID == 1 {
		SendRequest(&myPeer, verbose)

		//basta guardare il file peer_1.log
		<-myPeer.ChanStartTest
		executeTest(num)
	}
	if numSender == 2 && (myPeer.ID == 1 || myPeer.ID == 2) {
		SendRequest(&myPeer, verbose)
		<-myPeer.ChanStartTest
		//devo vedere entrambi i file dei sender
		time.Sleep(time.Minute / 2)
		if myPeer.ID == 1 {
			executeTest(num)
		}

	} else {
		time.Sleep(time.Minute / 2)
	}
}

func executeTest(num int) {
	numSender = num
	logPaths = list.New().Init()

	for i := 1; i < num+1; i++ {
		LogPath := "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(i) + ".log"
		logPaths.PushBack(LogPath)
	}

	fmt.Println(" ####################### TEST TOKEN-ASKING ", numSender, " sender #############################")

	testNoStarvation()
	if numSender == 2 {
		testSafety()
	}
	testMessageNumber()
}

func testMessageNumber() {

	index := 1
	for e := logPaths.Front(); e != nil; e = e.Next() {
		numMsg := 0
		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {

			if strings.Contains(fileScanner.Text(), "sends REQUEST message") ||
				strings.Contains(fileScanner.Text(), "receives TOKEN message") ||
				strings.Contains(fileScanner.Text(), "sends TOKEN message") {
				numMsg++
			}
		}

		if numMsg == 3 {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : PASSED !!!\n", index)
		} else {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : FAILED !!!\n", index)
		}
		index++
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
