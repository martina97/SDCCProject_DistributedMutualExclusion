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
	fmt.Println("sto in ExecuteTestPeer")
	myPeer = *peer

	if numSender == 1 && myPeer.ID == 1 {
		fmt.Println("mando il msg")
		SendRequest(&myPeer)
	}
	if numSender == 2 && (myPeer.ID == 1 || myPeer.ID == 2) {
		fmt.Println("mando il msg")
		SendRequest(&myPeer)
	} else {
		fmt.Println("sleep")
		time.Sleep(time.Minute / 2)
	}

}

func ExecuteTestCoordinator(coordinator *Coordinator, num int) {
	numSender = num
	logPaths = list.New().Init()
	//logPaths.Init()
	fmt.Println("sto in ExecuteTestCoordinator")

	myCoordinator = *coordinator

	//aspetta finche il numero di token msg ricevuti è pari a numSender
	//Wait connection
	for numMsg < numSender { //todo: mettere 3 , anche sotto
		ch := <-Connection
		if ch == true {
			numMsg++
		}
	}
	//fmt.Println("sto qua")
	Wg.Add(-numSender)
	//fmt.Println("sto qua2")

	for i := 0; i < utilities.MAXPEERS; i++ {
		if i != myCoordinator.ID {
			//fmt.Println(i)
			LogPath := "/docker/node_volume/tokenAsking/peer_" + strconv.Itoa(i) + ".log"
			logPaths.PushBack(LogPath)

			//fmt.Println(LogPath)
		}
	}
	testNoStarvation()
	if numSender == 2 {
		//fmt.Println("test safety !!!! ")
		testSafety()
	}
	testMessageNumber()

	/*
		ora posso controllare i vari file di log!!
		1 coordinator.log
		n-1 peer_n.log
	*/

	/*
		fileScanner := getFileSplit(myCoordinator.LogPath)
		for fileScanner.Scan() {
			//line := fileScanner.Text()

			fmt.Println(fileScanner.Text())
		}

	*/

	//f.Close()
}

func testMessageNumber() {

	for e := logPaths.Front(); e != nil; e = e.Next() {
		numMsg := 0
		fileScanner := utilities.GetFileSplit(e.Value.(string))
		for fileScanner.Scan() {
			//line := fileScanner.Text()

			//fmt.Println(fileScanner.Text())
			if strings.Contains(fileScanner.Text(), "send REQUEST message") ||
				strings.Contains(fileScanner.Text(), "receive TOKEN message") ||
				strings.Contains(fileScanner.Text(), "send TOKEN message") {
				//fmt.Println("CONTIENE !!!!! ")
				numMsg++
			}
		}
		fmt.Println("numMsg ===", numMsg)

		if numMsg == 3 {
			fmt.Println(" === TEST NUMBER OF MESSAGES: PASSED !!")
		} else {
			fmt.Println(" === TEST NUMBER OF MESSAGES: FAILED !!")
		}

		if numSender == 1 {
			break
		}
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

		fmt.Println("\n---------------------------------\n\n")
		index++

	}
	fmt.Println("enterP1 ==", enterP1)
	fmt.Println("exitP1 ==", exitP1)
	fmt.Println("enterP2 ==", enterP2)
	fmt.Println("exitP2 ==", exitP2)

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
