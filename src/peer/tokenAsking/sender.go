package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"net"
	"time"
)

var myPeer TokenPeer
var myCoordinator Coordinator

func SendRequest(peer *TokenPeer) {
	if myPeer.Username == "" { //vuol dire che non ho ancora inizializzato il peer
		myPeer = *peer
	}

	myPeer.mutex.Lock()
	utilities.WriteInfosToFile("tries to get the token.", myPeer.LogPath, myPeer.Username)
	//incremento Vector Clock!!!
	IncrementVC(myPeer.VC, myPeer.Username)

	date := time.Now().Format(utilities.DateFormat)
	msg := NewRequest(myPeer.Username, date, myPeer.VC)
	utilities.WriteVCInfoToFile(myPeer.LogPath, myPeer.Username, ToString(myPeer.VC))

	//mando REQUEST al coordinatore (è un campo di myPeer)
	connection := myPeer.Coordinator.Address + ":" + myPeer.Coordinator.Port
	addr, err := net.ResolveTCPAddr("tcp", connection)
	utilities.CheckError(err, "Unable to resolve IP")

	conn, err := net.DialTCP("tcp", nil, addr)
	err = conn.SetKeepAlive(true)
	utilities.CheckError(err, "Unable to set keepalive")

	enc := gob.NewEncoder(conn)
	err = enc.Encode(msg)
	utilities.CheckError(err, "Unable to encode message")

	msg.Receiver = utilities.COORDINATOR
	err = WriteMsgToFile("send", *msg, myPeer.LogPath, false)
	utilities.CheckError(err, "Error writing file")

	//invio msg di programma agli altri peer
	sendProgramMessage()

	myPeer.mutex.Unlock()

	go checkHasToken()

	//<-myPeer.HasToken

	myCoordinator.mutex.Unlock()
}

func sendProgramMessage() {
	date := time.Now().Format(utilities.DateFormat)
	msg := NewProgramMessage(myPeer.Username, date, myPeer.VC)

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		receiver := e.Value.(utilities.NodeInfo)
		if receiver.Username != utilities.COORDINATOR && receiver.Username != myPeer.Username {
			//open connection with peer
			peerConn := receiver.Address + ":" + receiver.Port
			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			utilities.CheckError(err, "Send response error on Dial")
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)
			msg.Receiver = receiver.Username
			err = WriteMsgToFile("send", *msg, myPeer.LogPath, false)
			utilities.CheckError(err, "error writing msg")
		}
	}
}

func sendToken(receiver string, isCoord bool) {

	if isCoord {
		for e := myCoordinator.PeerList.Front(); e != nil; e = e.Next() {
			dest := e.Value.(utilities.NodeInfo)
			if dest.Username == receiver {
				date := time.Now().Format(utilities.DateFormat)
				msg := NewTokenMessage(date, "coordinator", receiver, myCoordinator.VC)

				peerConn := dest.Address + ":" + dest.Port
				conn, err := net.Dial("tcp", peerConn)
				defer conn.Close()
				utilities.CheckError(err, "Send response error on Dial")

				enc := gob.NewEncoder(conn)
				enc.Encode(msg)
				err = WriteMsgToFile("send", *msg, myCoordinator.LogPath, true)
				utilities.CheckError(err, "error writing msg")
			}
		}
	} else {
		date := time.Now().Format(utilities.DateFormat)
		msg := NewTokenMessage(date, myPeer.Username, "coordinator", myPeer.VC)
		connection := myPeer.Coordinator.Address + ":" + myPeer.Coordinator.Port

		conn, err := net.Dial("tcp", connection)
		defer conn.Close()
		utilities.CheckError(err, "Send response error on Dial")

		enc := gob.NewEncoder(conn)
		enc.Encode(msg)
		err = WriteMsgToFile("send", *msg, myPeer.LogPath, false)
		utilities.CheckError(err, "error writing msg")
	}
}
