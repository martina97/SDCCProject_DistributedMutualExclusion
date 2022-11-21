package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

var myPeer TokenPeer
var myCoordinator Coordinator

func SendRequest(peer *TokenPeer) {
	if myPeer.Username == "" { //vuol dire che non ho ancora inizializzato il peer
		fmt.Println("sto in SendRequest --- RA_PEER VUOTA")
		myPeer = *peer

	} else {
		fmt.Println("sto in SendRequest --- RA_PEER NON VUOTA")
	}
	fmt.Println("myPeer.Coordinator= ", myPeer.Coordinator)

	myPeer.mutex.Lock()
	//incremento Vector Clock!!!
	fmt.Println("myTokenPeer.VC =", myPeer.VC)
	fmt.Println("incremento VC")
	utilities.IncrementVC(myPeer.VC, myPeer.Username)

	fmt.Println("myTokenPeer.VC =", myPeer.VC)
	date := time.Now().Format(utilities.DATE_FORMAT)
	msg := NewRequest(myPeer.Username, date, myPeer.VC)
	fmt.Println("IL MESSAGGIO E' ====", msg.ToString("send"))

	utilities.WriteVCInfoToFile(myPeer.LogPath, myPeer.Username, myPeer.VC)
	fmt.Println("dopo WriteVCInfoToFile")

	//ora mando REQUEST al coordinatore (è un campo di myPeer)
	fmt.Println("devo inviare req al coordinatore")
	fmt.Println(" myPeer.Coordinator.Address  ==", myPeer.Coordinator.Address)
	fmt.Println(" myPeer.Coordinator.Port  ==", myPeer.Coordinator.Port)
	connection := myPeer.Coordinator.Address + ":" + myPeer.Coordinator.Port

	conn, err := net.Dial("tcp", connection)
	fmt.Println("dopo Dial")
	defer conn.Close()
	if err != nil {
		log.Println("Send response error on Dial")
	}
	enc := gob.NewEncoder(conn)
	fmt.Println("dopo NewEncoder")
	enc.Encode(msg)

	fmt.Println("dopo encode msg, msg == ", msg.ToString("send"))
	msg.Receiver = utilities.COORDINATOR
	fmt.Println("dopo encode msg, msg == ", msg.ToString("send"))
	err = WriteMsgToFile("send", *msg)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}

	myPeer.mutex.Unlock()
}
