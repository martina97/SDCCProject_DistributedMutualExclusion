package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"net"
	"strconv"
	"time"
)

var myPeer TokenPeer

func SendRequest(peer *TokenPeer) {
	if myPeer.Username == "" { //vuol dire che non ho ancora inizializzato il peer
		myPeer = *peer
	}

	utilities.SleepRandInt()

	myPeer.mutex.Lock()
	utilities.WriteInfosToFile("tries to get the token.", myPeer.LogPath, myPeer.Username)
	//incremento Vector Clock!!!
	IncrementVC(myPeer.VC, myPeer.Username)

	date := time.Now().Format(utilities.DateFormat)
	msg := NewRequest(myPeer.Username, date, myPeer.VC)
	utilities.WriteVCInfoToFile(myPeer.LogPath, myPeer.Username, ToString(myPeer.VC))

	connection := utilities.CoordAddr + ":" + strconv.Itoa(utilities.ServerPort)

	addr, err := net.ResolveTCPAddr("tcp", connection)
	utilities.CheckError(err, "Unable to resolve IP")

	conn, err := net.DialTCP("tcp", nil, addr)
	err = conn.SetKeepAlive(true)
	utilities.CheckError(err, "Unable to set keepalive")

	enc := gob.NewEncoder(conn)
	err = enc.Encode(msg)
	utilities.CheckError(err, "Unable to encode message")

	msg.Receiver = "coordinator"
	err = WriteMsgToFile("send", *msg, myPeer.LogPath, false)
	utilities.CheckError(err, "Error writing file")

	//invio msg di programma agli altri peer
	sendProgramMessage()

	myPeer.mutex.Unlock()

	go checkHasToken()

}

func sendProgramMessage() {
	date := time.Now().Format(utilities.DateFormat)
	msg := NewProgramMessage(myPeer.Username, date, myPeer.VC)

	for e := myPeer.PeerList.Front(); e != nil; e = e.Next() {
		receiver := e.Value.(utilities.NodeInfo)
		if receiver.Username != "coordinator" && receiver.Username != myPeer.Username {
			utilities.SleepRandInt()

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

func sendToken() {

	utilities.SleepRandInt()
	date := time.Now().Format(utilities.DateFormat)
	msg := NewTokenMessage(date, myPeer.Username, "coordinator", myPeer.VC)
	connection := utilities.CoordAddr + ":" + strconv.Itoa(utilities.ServerPort)

	conn, err := net.Dial("tcp", connection)
	defer conn.Close()
	utilities.CheckError(err, "Send response error on Dial")

	enc := gob.NewEncoder(conn)
	enc.Encode(msg)
	err = WriteMsgToFile("send", *msg, myPeer.LogPath, false)
	myPeer.ChanStartTest <- true
	utilities.CheckError(err, "error writing msg")

}
