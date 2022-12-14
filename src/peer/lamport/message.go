package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
	"log"
	"strconv"
	"time"
)

type MessageType string

const (
	Request MessageType = "Request"
	Reply   MessageType = "Reply"
	Release MessageType = "Release"
)

type Message struct {
	MsgType  MessageType //request,reply,release
	Sender   string
	Receiver string
	Date     string
	TS       ScalarClock
}

func NewRequest(sender string, date string, timeStamp ScalarClock) *Message {
	return &Message{
		MsgType: Request,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perché la richiesta viene mandata a tutti
		Date: date,
		TS:   timeStamp,
	}
}

func NewReply(sender string, receiver string, date string, timeStamp ScalarClock) *Message {
	return &Message{
		MsgType:  Reply,
		Sender:   sender,
		Receiver: receiver,
		Date:     date,
		TS:       timeStamp,
	}
}

func NewRelease(sender string, date string, timeStamp ScalarClock) *Message {
	return &Message{

		MsgType: Release,
		Sender:  sender,
		Date:    date,
		TS:      timeStamp,
	}
}

func (m *Message) ToString(role string) string {
	var name string

	switch m.MsgType {
	case Request:
		name = "Request"
	case Release:
		name = "Release"
	case Reply:
		name = "Reply"
	}

	if role == "send" {
		return fmt.Sprintf(" %s message: {%s %s %s [%d]}", name, name, m.Receiver, m.Date, m.TS)
	}
	if role == "receive" {
		return fmt.Sprintf(" %s message: {%s %s [%d]} from %s", name, name, m.Date, m.TS, m.Sender)
	}

	return ""
}

func WriteMsgToFile(path string, id string, typeMsg string, message Message, timestamp ScalarClock) {

	var err error
	f := utilities.OpenFile(path)
	//save new address on file
	date := time.Now().Format(utilities.DateFormat)
	if typeMsg == "send" {
		_, err = f.WriteString("[" + date + "] : " + id + " sends" + message.ToString("send") + " to " + message.Receiver + ".")
	}
	if typeMsg == "receive" {
		_, err = f.WriteString("[" + date + "] : " + id + " receives" + message.ToString("receive"))
		if message.MsgType != Reply {
			_, err = f.WriteString(" and updates its logical scalar clock to " + strconv.Itoa(int(timestamp)))
		}
	}
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}
