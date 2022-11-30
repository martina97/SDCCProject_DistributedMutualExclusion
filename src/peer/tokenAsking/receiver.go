package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"time"
)

var verbose bool

func HandleConnectionPeer(conn net.Conn, peer *TokenPeer) error {

	flag.BoolVar(&verbose, "v", utilities.Verbose, "use this flag to get verbose info on messages")
	flag.Parse()

	if verbose {
		fmt.Println("VERBOSE FLAG ON")
	}

	if myPeer.Username == "" {
		myPeer = *peer
	}
	defer conn.Close()

	msg := new(Message)
	dec := gob.NewDecoder(conn)
	err := dec.Decode(msg)
	utilities.CheckError(err, "error decoding msg")

	if msg.MsgType == ProgramMessage {
		myPeer.mutex.Lock()
		//update VC !
		UpdateVC(myPeer.VC, msg.VC)
		if verbose {
			WriteMsgToFile("receive", *msg, myPeer.LogPath, false)
		}
		utilities.CheckError(err, "error writing msg")
		myPeer.mutex.Unlock()
	}

	if msg.MsgType == Token {
		// ho il token !
		myPeer.mutex.Lock()
		if verbose {
			WriteMsgToFile("receive", *msg, myPeer.LogPath, false)
		}
		utilities.CheckError(err, "error writing msg")
		myPeer.HasToken = true
		myPeer.mutex.Unlock()
	}

	return nil
}

func checkHasToken() {

	for !(myPeer.HasToken) {

		time.Sleep(time.Second * 5)
	}
	myPeer.mutex.Lock()
	date := time.Now().Format(utilities.DateFormat)
	if verbose {
		utilities.WriteInfosToFile("enters the critical section at "+date+".", myPeer.LogPath, myPeer.Username)
	}
	time.Sleep(time.Second * 15)
	date = time.Now().Format(utilities.DateFormat)

	if verbose {
		utilities.WriteInfosToFile("exits the critical section at "+date+".", myPeer.LogPath, myPeer.Username)
		utilities.WriteInfosToFile("releases the token.", myPeer.LogPath, myPeer.Username)
	}

	//invio msg con il token al coordinatore
	sendToken()
	myPeer.mutex.Unlock()

}
