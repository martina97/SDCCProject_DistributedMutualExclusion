package RicartAgrawala

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
)

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

//Save message
func HandleConnection(conn net.Conn) error {

	fmt.Println("sto in handleConnection dentro RicartAgrawala package")
	if MyRApeer == (RApeer{}) {
		fmt.Println("RA_PEER VUOTA")
	} else {
		fmt.Println("RA_PEER NON VUOTA")
	}

	//devo vedere se già ho inizializzato RApeer (se non ho mai inviato un msg, non l'ho inizializzato)
	fmt.Println("MyRApeer == ", MyRApeer.ToString())

	defer conn.Close()
	msg := new(utilities.Message)

	dec := gob.NewDecoder(conn)
	dec.Decode(msg)

	mutex := MyRApeer.GetMutex()
	if msg.MsgType == utilities.Request {
		fmt.Println("MESS REQUEST !!!!!! ")
		mutex.Lock()
		utilities.UpdateTS(&MyRApeer.Num, &msg.TS, "RicartAgrawala")

		utilities.WriteMsgToFile3(MyRApeer.LogPath, MyRApeer.Username, "Receive", *msg, MyRApeer.Num, "RicartAgrawala")

	}

	return nil
}
