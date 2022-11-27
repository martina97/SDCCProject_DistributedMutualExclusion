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
		<-MyRApeer.ChanAcquireLock
	}
	if numSender == 2 && (MyRApeer.ID == 0 || MyRApeer.ID == 1) {
		fmt.Println("mando il msg")
		SendRicart(&MyRApeer)
		<-MyRApeer.ChanAcquireLock

	} else {
		fmt.Println("sleep")
		time.Sleep(time.Minute)
	}

	fmt.Println(" ####################### TEST #############################")

	for i := 0; i < num; i++ {
		//fmt.Println(i)
		LogPath := "/docker/node_volume/ricartAgrawala/peer_" + strconv.Itoa(i) + ".log"
		logPaths.PushBack(LogPath)
	}

	//faccio eseguire a p1 i test, ossia legge tutti i file (perchè è l'ultimo che esegue
	// l'algoritmo)
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
		2(N-1) messaggi per accedere alla CS:
			• N-1 messaggi di richiesta
			• N-1 messaggi di reply
	*/

	index := 0
	for e := logPaths.Front(); e != nil; e = e.Next() {
		numMsg := 0
		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {
			//line := fileScanner.Text()

			//fmt.Println(fileScanner.Text())
			if strings.Contains(fileScanner.Text(), "sends Request message") ||
				strings.Contains(fileScanner.Text(), "receives Reply message") {
				//fmt.Println("CONTIENE !!!!! ")
				numMsg++
			}
		}
		//fmt.Println("numMsg ===", numMsg)

		if numMsg == 2*(utilities.MAXPEERS-1) {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : PASSED !!!\n", index)
		} else {
			fmt.Printf(" === TEST NUMBER OF MESSAGES p%d : FAILED !!!\n", index)
		}
		index++

	}
}
