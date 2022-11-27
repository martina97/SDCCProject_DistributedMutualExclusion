package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

type MessageType int

type TiS []uint64

const (
	Request MessageType = iota + 1 // request mutual lock
	Reply                          // reply mutual lock
	Release                        // release mutual lock

)

/*
servono lettera maiuscola perche devo farne il marshaling e unmarshaling, che richiedono solo field esportabili
todo: vedi libro pag 108
*/
type Message struct {

	//MsgID   string
	MsgType MessageType //request,reply,release

	//MsgContent interface{}
	Sender     string
	SenderProc utilities.NodeInfo
	Receiver   string
	SeqNum     []uint64
	Date       string
	TS         utilities.TimeStamp
}

// NewRequest returns a new distributed mutual lock message.
func NewRequest(sender string, date string, timeStamp utilities.TimeStamp) *Message {
	return &Message{
		MsgType: Request,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		Date: date,
		TS:   timeStamp,
	}
}

func NewReply(sender string, receiver string, date string, timeStamp utilities.TimeStamp) *Message {
	return &Message{
		MsgType:  Reply,
		Sender:   sender,
		Receiver: receiver,
		Date:     date,
		TS:       timeStamp,
	}
}
func NewRelease(sender string, date string, timeStamp utilities.TimeStamp) *Message {
	return &Message{

		MsgType: Release,
		Sender:  sender,
		Date:    date,
		TS:      timeStamp,
	}
}

/*
IL MESSAGGIO E' ==== {prova 1 2 [1] 2022/04/04 10:37:35}
ID MESSAGGIO E' ==== prova
MsgType MESSAGGIO E' ==== 1 (Request)
Sender MESSAGGIO E' ==== 2
SeqNum MESSAGGIO E' ==== [1]
Date MESSAGGIO E' ==== 2022/04/04 10:37:35

*/

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
		//Request message: {Request [] p3 p1 17:39:42.230 [1]} to p0.
		//return fmt.Sprintf(" %s message: {%s %s %s %s [%d]}", name, name, m.Sender, m.Receiver, m.Date, m.TS)

		//Request message: {Request [] p1 17:39:42.230 [1]} --- p1=receiver, [1] = timestamp
		return fmt.Sprintf(" %s message: {%s %s %s [%d]}", name, name, m.Receiver, m.Date, m.TS)
	}
	if role == "receive" {
		return fmt.Sprintf(" %s message: {%s %s [%d]} from %s", name, name, m.Date, m.TS, m.Sender)
	}

	return ""
}
