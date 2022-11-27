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
	VC       VectorClock
}

func NewRequest(sender string, date string, vc VectorClock) *Message {
	return &Message{
		MsgType: Request,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata solo al coordinatore
		Date: date,
		VC:   vc,
	}
}

func NewProgramMessage(sender string, date string, vc VectorClock) *Message {
	return &Message{
		MsgType: ProgramMessage,
		Sender:  sender,
		Date:    date,
		VC:      vc,
	}
}
func NewTokenMessage(date string, sender string, receiver string, vc VectorClock) *Message {
	return &Message{
		MsgType:  Token,
		Sender:   sender,
		Receiver: receiver,
		Date:     date,
		VC:       vc,
	}
}

func (m *Message) ToString(role string) string {
	var name string

	switch m.MsgType {
	case Request:
		name = "REQUEST"
	case ProgramMessage:
		name = "PROGRAM"
	case Token:
		name = "TOKEN"
	}

	if role == "send" {
		//Request message: {Request [] p3 p1 17:39:42.230 [1]} to p0.
		//return fmt.Sprintf(" %s message: {%s %s %s %s [%d]}", name, name, m.Sender, m.Receiver, m.Date, m.TS)

		//Request message: {Request [] p1 17:39:42.230 [1]} --- p1=receiver, [1] = timestamp
		return fmt.Sprintf(" %s message: {%s %s %s %s}", name, m.MsgType, m.Receiver, m.Date, ToString(m.VC))
	}
	if role == "receive" {
		return fmt.Sprintf(" %s message: {%s %s %s} from %s", name, m.MsgType, m.Date, ToString(m.VC), m.Sender)
	}

	return ""
}

func WriteMsgToFile(action string, message Message, isCoord bool) error {

	var username string
	var err error

	f := openFile(isCoord)
	if isCoord {
		username = "coordinator"
	} else {
		username = myPeer.Username
	}

	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)

	if action == "send" {
		switch message.MsgType {
		case Request:
			_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " sends" + message.ToString("send") + " to coordinator.")
		case ProgramMessage:
			_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " sends" + message.ToString("send") + " to " + message.Receiver + ".")
		case Token:
			if isCoord {
				_, err = f.WriteString("[" + date + "] : coordinator sends" + message.ToString("send") + " to " + message.Receiver + ".")
			} else {
				_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " sends" + message.ToString("send") + " to coordinator.")
			}

		}
	} else {
		switch message.MsgType {
		case ProgramMessage:
			_, err = f.WriteString("[" + date + "] : " + myPeer.Username + " receives" + message.ToString("receive") + " and update its vector clock to " + ToString(myPeer.VC) + ".")
		case Request:
			_, err = f.WriteString("[" + date + "] : coordinator receives" + message.ToString("receive") + ".")
		case Token:
			_, err = f.WriteString("[" + date + "] : " + username + " receives" + message.ToString("receive") + ".\n")
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

func WriteInfosToFile(text string, isCoord bool) {
	var username string

	f := openFile(isCoord)
	if isCoord {
		username = "coordinator"
	} else {
		username = myPeer.Username
	}

	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)

	_, _ = f.WriteString("[" + date + "] : " + username + " " + text)

	_, _ = f.WriteString("\n")
	_ = f.Sync()
}

func openFile(isCoord bool) *os.File {
	var err error
	var f *os.File

	if isCoord {
		f, err = os.OpenFile(myCoordinator.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	} else {
		f, err = os.OpenFile(myPeer.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	}
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}

func WriteVCInfoToFile(isCoord bool) {
	var vc VectorClock
	var username string

	f := openFile(isCoord)

	if isCoord {
		vc = myCoordinator.VC
		username = "coordinator"
	} else {
		vc = myPeer.VC
		username = myPeer.Username
	}

	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)

	_, err := f.WriteString("[" + date + "] : " + username + " update its vector clock to " + ToString(vc))
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}
