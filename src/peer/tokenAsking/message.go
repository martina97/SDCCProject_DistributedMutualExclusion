package tokenAsking

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

type MessageType string

const (
	Request        MessageType = "REQUEST"
	ProgramMessage MessageType = "PROGRAM_MESSAGE"
	Token          MessageType = "TOKEN"
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

func (m *Message) ToString(role string) string {
	var name string
	//date := time.Now().Format("2006/01/02 15:04:05")

	switch m.MsgType {
	case Request:
		name = "Request"
	case ProgramMessage:
		name = "ProgramMessage"
	case Token:
		name = "Token"
	}

	fmt.Println("sto in ToString -----", m.Sender)
	fmt.Println("sto in ToString -----", m.Receiver)
	if role == "send" {
		//Request message: {Request [] p3 p1 17:39:42.230 [1]} to p0.
		//return fmt.Sprintf(" %s message: {%s %s %s %s [%d]}", name, name, m.Sender, m.Receiver, m.Date, m.TS)

		//Request message: {Request [] p1 17:39:42.230 [1]} --- p1=receiver, [1] = timestamp
		return fmt.Sprintf(" %s message: {%s %s %s [%d]}", name, name, m.Receiver, m.Date, utilities.ToString(m.VC))
	}
	if role == "receive" {
		return fmt.Sprintf(" %s message: {%s %s [%d]} from %s", name, name, m.Date, utilities.ToString(m.VC), m.Sender)
	}

	return ""
}
