package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
	"log"
	"os"
	"time"
)

type MessageType string

const (
	Request        MessageType = "Request"
	ProgramMessage MessageType = "Program message"
	Token          MessageType = "Token"
)

type Message struct {
	MsgType  MessageType //Request,ProgramMessage,Token
	Sender   string
	Receiver string
	Date     string
	VC       utilities.VectorClock
}

func NewRequest(sender string, date string, vc utilities.VectorClock) *Message {
	return &Message{
		//MsgID:      RandStringBytes(msgIDCnt),
		MsgType: Request,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		Date: date,
		VC:   vc,
	}
}

func NewProgramMessage(sender string, date string, vc utilities.VectorClock) *Message {
	return &Message{
		MsgType: ProgramMessage,
		Sender:  sender,
		Date:    date,
		VC:      vc,
	}
}
func NewTokenMessage(date string, receiver string, vc utilities.VectorClock) *Message {
	return &Message{
		MsgType:  Token,
		Sender:   "coordinator",
		Receiver: receiver,
		Date:     date,
		VC:       vc,
	}
}

func (m *Message) ToString(role string) string {
	var name string
	//date := time.Now().Format("2006/01/02 15:04:05")

	switch m.MsgType {
	case Request:
		name = "REQUEST"
	case ProgramMessage:
		name = "PROGRAM"
	case Token:
		name = "TOKEN"
	}

	fmt.Println("sto in ToString -----", m.Sender)
	fmt.Println("sto in ToString -----", m.Receiver)
	if role == "send" {
		//Request message: {Request [] p3 p1 17:39:42.230 [1]} to p0.
		//return fmt.Sprintf(" %s message: {%s %s %s %s [%d]}", name, name, m.Sender, m.Receiver, m.Date, m.TS)

		//Request message: {Request [] p1 17:39:42.230 [1]} --- p1=receiver, [1] = timestamp
		return fmt.Sprintf(" %s message: {%s %s %s %s}", name, m.MsgType, m.Receiver, m.Date, utilities.ToString(m.VC))
	}
	if role == "receive" {
		return fmt.Sprintf(" %s message: {%s %s %s} from %s", name, m.MsgType, m.Date, utilities.ToString(m.VC), m.Sender)
	}

	return ""
}

func WriteMsgToFile(action string, message Message, isCoord bool) error {
	fmt.Println("sto in WriteMsgToFile")
	fmt.Println("path == ", myPeer.LogPath)
	var f *os.File
	var err error
	var username string
	if isCoord {
		f, err = os.OpenFile(myCoordinator.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		username = "coordinator"
	} else {
		f, err = os.OpenFile(myPeer.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		username = myPeer.Username
	}
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)

	if action == "send" {
		switch message.MsgType {
		case Request:
			_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " " + action + message.ToString("send") + " to coordinator.")
		case ProgramMessage:
			_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " " + action + message.ToString("send") + " to " + message.Receiver + ".")
		case Token:
			if isCoord {
				_, err = f.WriteString("[" + date + "] : coordinator " + action + message.ToString("send") + " to " + message.Receiver + ".")
			} else {
				_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " " + action + message.ToString("send") + " to coordinator.")
			}

		}
	} else {
		switch message.MsgType {
		case ProgramMessage:
			_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " " + action + message.ToString("receive") + " and update its vector clock to " + utilities.ToString(myPeer.VC) + ".")
		case Request:
			_, err = f.WriteString("[" + date + "] : coordinator " + action + message.ToString("receive") + ".")
		case Token:
			_, err = f.WriteString("[" + date + "] : " + username + " " + action + message.ToString("receive") + ".\n")
			_, err = f.WriteString("[" + date + "] : " + username + " gets the token.")
		}

	}
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil

}

func WriteInfosToFile(text string) {
	f, err := os.OpenFile(myPeer.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)

	_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " " + text)

	_, err = f.WriteString("\n")
	err = f.Sync()
}

func WriteVCInfoToFile(isCoord bool) {
	var f *os.File
	var err error

	var vc utilities.VectorClock
	var username string

	if isCoord {
		f, err = os.OpenFile(myCoordinator.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		vc = myCoordinator.VC
		username = "coordinator"

	} else {
		f, err = os.OpenFile(myPeer.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		vc = myPeer.VC
		username = myPeer.Username

	}
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)

	_, err = f.WriteString("[" + date + "] : " + username + " update its vector clock to " + utilities.ToString(vc))
	_, err = f.WriteString("\n")
	err = f.Sync()
}
