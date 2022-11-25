package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"container/list"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
	numSender  int
	logPaths   *list.List
	numMsg     int
)

func ExecuteTestPeer(peer *LamportPeer, num int) {
	numSender = num
	logPaths = list.New().Init()

	fmt.Println("sto in ExecuteTestPeer")
	myPeer = *peer

	if numSender == 1 && myPeer.ID == 0 {
		fmt.Println("mando il msg")
		SendLamport(&myPeer)
	}
	if numSender == 2 && (myPeer.ID == 0 || myPeer.ID == 1) {
		fmt.Println("mando il msg")
		SendLamport(&myPeer)
	} else {
		fmt.Println("sleep")
		time.Sleep(time.Minute + time.Minute/2)
	}

	fmt.Println(" ####################### TEST #############################")

	for i := 0; i < num; i++ {
		//fmt.Println(i)
		LogPath := "/docker/node_volume/ricartAgrawala/peer_" + strconv.Itoa(i) + ".log"
		logPaths.PushBack(LogPath)
	}

	//faccio eseguire a p2 i test, ossia legge tutti i file, per farlo devo aspettare che riceva numSender msg di release!!
	if myPeer.ID == 2 {

		fmt.Println("numMsg ==", numMsg)
		fmt.Println("numSender ==", numSender)
		fmt.Println("wg ==", Wg)
		for numMsg < numSender {
			ch := <-Connection
			if ch == true {
				numMsg++
			}
		}
		//fmt.Println("sto qua")
		Wg.Add(-numSender)

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
			//line := fileScanner.Text()

			//fmt.Println(fileScanner.Text())
			if strings.Contains(fileScanner.Text(), "enters the critical section") {
				//fmt.Println("CONTIENE !!!!! ")
				csEntries++
			}
		}
		/*
			if numSender == 1 {
				break
			}

		*/
		//fmt.Println("\n---------------------------------\n\n")
	}
	//fmt.Println("csEntries == ", csEntries)

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
		//var index int //mi serve per vedere a quante iterazioni sto

		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {
			//line := fileScanner.Text()
			//fmt.Println(fileScanner.Text())
			if strings.Contains(fileScanner.Text(), stringEnter) {
				//fmt.Println("CONTIENE !!!!! ")
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

		//fmt.Println("\n---------------------------------\n\n")
		index++

	}
	/*
		fmt.Println("enterP1 ==", enterP1)
		fmt.Println("exitP1 ==", exitP1)
		fmt.Println("enterP2 ==", enterP2)
		fmt.Println("exitP2 ==", exitP2)

	*/

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
		3(N-1) messaggi per accedere alla CS:
			• N-1 messaggi di richiesta
			• N-1 messaggi di reply
			• N-1 messaggi di release
	*/

	index := 0
	for e := logPaths.Front(); e != nil; e = e.Next() {
		numMsg := 0
		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {
			//line := fileScanner.Text()

			fmt.Println(fileScanner.Text())
			if strings.Contains(fileScanner.Text(), "send Request message") ||
				strings.Contains(fileScanner.Text(), "receive Reply message") ||
				strings.Contains(fileScanner.Text(), "send Release message") {
				//fmt.Println("CONTIENE !!!!! ")
				numMsg++
			}
		}
		//fmt.Println("numMsg ===", numMsg)

		if numMsg == 3*(utilities.MAXPEERS-1) {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : PASSED !!!\n", index)
		} else {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : FAILED !!!\n", index)
		}
		index++

	}
}
