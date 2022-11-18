package RicartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"net"
)

/*
func Message_handler() {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(utilities.Client_port))
	if err != nil {
		log.Fatal("net.Lister fail")
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept fail")
		}
		go HandleConnection(connection)
		//go handleConnectionCentralized(connection)
	}
}

*/

//Save message
func HandleConnection(conn net.Conn, peer *RApeer) error {

	fmt.Println("sto in handleConnection dentro RicartAgrawala package")
	if MyRApeer == (RApeer{}) {
		fmt.Println("RA_PEER VUOTA")
		MyRApeer = *peer
	} else {
		fmt.Println("RA_PEER NON VUOTA")
	}

	//devo vedere se già ho inizializzato RApeer (se non ho mai inviato un msg, non l'ho inizializzato)
	fmt.Println("MyRApeer == ", MyRApeer.ToString())

	defer conn.Close()
	msg := new(utilities.Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)
	fmt.Println("il msg == ", msg.MessageToString("receive"))

	mutex := MyRApeer.GetMutex()
	if msg.MsgType == utilities.Request {

		/*
			Upon receipt REQUEST(t) from pj
			1. if State=CS or (State=Requesting and {Last_Req, i} < {t, j})
				then insert {t, j} in Q
			3. else
				send REPLY to pj
			4. Num = max(t, Num)
		*/
		fmt.Println("MESS REQUEST !!!!!! ")
		mutex.Lock()
		utilities.UpdateTS(&MyRApeer.Num, &msg.TS, "RicartAgrawala")

		utilities.WriteMsgToFile3(MyRApeer.LogPath, MyRApeer.Username, "Receive", *msg, MyRApeer.Num, "RicartAgrawala")

		checkConditions(msg)
		mutex.Unlock()

	}

	return nil
}

func checkConditions(msg *utilities.Message) bool {

	if (MyRApeer.state == CS) || (MyRApeer.state == Requesting && checkTS(msg)) {
		fmt.Println("sto in checkConditions -->  non invio reply e metto msg in coda")
		MyRApeer.DeferSet.PushBack(msg)
		return true
	}
	fmt.Println("invio reply!!!!!!")
	return false

}

func checkTS(msg *utilities.Message) bool {
	// true se {Last_Req, i} < {t, j})
	if (MyRApeer.lastReq < msg.TS) && (MyRApeer.Username < msg.Sender) {
		fmt.Println("sto in checkTS e la condizione e' true --> non invio reply e metto msg in coda")
		return true
	}
	return false

}
