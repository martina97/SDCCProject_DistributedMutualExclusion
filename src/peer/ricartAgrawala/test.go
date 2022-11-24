package ricartAgrawala

import (
	"container/list"
	"fmt"
	"time"
)

var numMsg int
var logPaths *list.List
var numSender int

func ExecuteTestPeer(peer *RApeer, num int) {
	numSender = num
	fmt.Println("sto in ExecuteTestPeer")
	MyRApeer = *peer

	if numSender == 1 && MyRApeer.ID == 0 {
		fmt.Println("mando il msg")
		SendRicart(&MyRApeer)
	}
	if numSender == 2 && (MyRApeer.ID == 1 || MyRApeer.ID == 2) {
		fmt.Println("mando il msg")
		SendRicart(&MyRApeer)
	} else {
		fmt.Println("sleep")
		time.Sleep(time.Minute / 2)
	}

}