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
	WriteInfosToFile("try to get the token.", false)
	//incremento Vector Clock!!!
	fmt.Println("myTokenPeer.VC =", myPeer.VC)
	fmt.Println("incremento VC")
	utilities.IncrementVC(myPeer.VC, myPeer.Username)

	fmt.Println("myTokenPeer.VC =", myPeer.VC)
	date := time.Now().Format(utilities.DATE_FORMAT)
	msg := NewRequest(myPeer.Username, date, myPeer.VC)
	fmt.Println("IL MESSAGGIO E' ====", msg.ToString("send"))

	WriteVCInfoToFile(false)
	fmt.Println("dopo WriteVCInfoToFile")

	//ora mando REQUEST al coordinatore (è un campo di myPeer)
	fmt.Println("devo inviare req al coordinatore")
	fmt.Println(" myPeer.Coordinator.Address  ==", myPeer.Coordinator.Address)
	fmt.Println(" myPeer.Coordinator.Port  ==", myPeer.Coordinator.Port)
	connection := myPeer.Coordinator.Address + ":" + myPeer.Coordinator.Port
	addr, err := net.ResolveTCPAddr("tcp", connection)
	if err != nil {
		fmt.Printf("Unable to resolve IP")
	}

	//conn, err := net.Dial("tcp", connection)
	conn, err := net.DialTCP("tcp", nil, addr)
	fmt.Println("dopo Dial")
	err = conn.SetKeepAlive(true)
	if err != nil {
		fmt.Printf("Unable to set keepalive - %s", err)
	}

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
	err = WriteMsgToFile("send", *msg, false)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}

	//ora invio msg di programma agli altri peer
	sendProgramMessage()

	myPeer.mutex.Unlock()
	<-myPeer.HasToken
	fmt.Println("ho il token!!!!")
	date = time.Now().Format(utilities.DATE_FORMAT)
	WriteInfosToFile("enters the critical section at "+date+".", false)
	time.Sleep(time.Minute / 2)
	fmt.Println("rilascio il token!!!!")
	date = time.Now().Format(utilities.DATE_FORMAT)
	WriteInfosToFile("exits the critical section at "+date+".", false)
	WriteInfosToFile("releases the token.", false)

	myCoordinator.mutex.Lock()
	//devo inviare msg con il token al coordinatore
	sendToken("coordinator", false)
	myCoordinator.mutex.Unlock()

}

func sendProgramMessage() {
	fmt.Println("sto in sendProgramMessage")
	date := time.Now().Format(utilities.DATE_FORMAT)
	msg := NewProgramMessage(myPeer.Username, date, myPeer.VC)

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		fmt.Println("sto nel for")
		fmt.Println("il msg == ", msg)
		receiver := e.Value.(utilities.NodeInfo)
		if receiver.Username != utilities.COORDINATOR && receiver.Username != myPeer.Username {
			fmt.Println("sto nell'if ")
			//open connection whit peer

			/*
				peerConn := receiver.Address + ":" + receiver.Port
				conn, err := net.Dial("tcp", peerConn)
				defer conn.Close()
				if err != nil {
					log.Println("Send response error on Dial")
				}
			*/
			connection := receiver.Address + ":" + receiver.Port
			addr, err := net.ResolveTCPAddr("tcp", connection)
			if err != nil {
				fmt.Printf("Unable to resolve IP")
			}

			//conn, err := net.Dial("tcp", connection)
			conn, err := net.DialTCP("tcp", nil, addr)
			fmt.Println("dopo Dial")
			err = conn.SetKeepAlive(true)
			if err != nil {
				fmt.Printf("Unable to set keepalive - %s", err)
			}
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			msg.Receiver = receiver.Username

			err = WriteMsgToFile("send", *msg, false)

			if err != nil {
				log.Fatalf("error writing msg %v", err)
			}
		}

	}
}

func sendToken(receiver string, isCoord bool) {
	fmt.Println("sto in sendToken")

	if isCoord {
		for e := myCoordinator.PeerList.Front(); e != nil; e = e.Next() {
			dest := e.Value.(utilities.NodeInfo)
			if dest.Username == receiver {
				date := time.Now().Format(utilities.DATE_FORMAT)
				msg := NewTokenMessage(date, "coordinator", receiver, myCoordinator.VC)

				/*
					peerConn := dest.Address + ":" + dest.Port
					conn, err := net.Dial("tcp", peerConn)
					defer conn.Close()
					if err != nil {
						log.Println("Send response error on Dial")
					}

				*/
				connection := dest.Address + ":" + dest.Port
				addr, err := net.ResolveTCPAddr("tcp", connection)
				if err != nil {
					fmt.Printf("Unable to resolve IP")
				}

				//conn, err := net.Dial("tcp", connection)
				conn, err := net.DialTCP("tcp", nil, addr)
				fmt.Println("dopo Dial")
				err = conn.SetKeepAlive(true)
				if err != nil {
					fmt.Printf("Unable to set keepalive - %s", err)
				}
				enc := gob.NewEncoder(conn)
				enc.Encode(msg)
				err = WriteMsgToFile("send", *msg, true)
				if err != nil {
					log.Fatalf("error writing msg %v", err)
				}
			}
		}
	} else {
		date := time.Now().Format(utilities.DATE_FORMAT)
		msg := NewTokenMessage(date, myPeer.Username, "coordinator", myPeer.VC)
		connection := myPeer.Coordinator.Address + ":" + myPeer.Coordinator.Port

		/*
			conn, err := net.Dial("tcp", coordConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}

		*/
		addr, err := net.ResolveTCPAddr("tcp", connection)
		if err != nil {
			fmt.Printf("Unable to resolve IP")
		}

		//conn, err := net.Dial("tcp", connection)
		conn, err := net.DialTCP("tcp", nil, addr)
		fmt.Println("dopo Dial")
		err = conn.SetKeepAlive(true)
		if err != nil {
			fmt.Printf("Unable to set keepalive - %s", err)
		}

		enc := gob.NewEncoder(conn)
		enc.Encode(msg)
		err = WriteMsgToFile("send", *msg, false)
		if err != nil {
			log.Fatalf("error writing msg %v", err)
		}
	}

}
