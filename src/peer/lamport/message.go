package lamport

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
	"log"
	"strconv"
	"time"
)

type MessageType int

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
	MsgType  MessageType //request,reply,release
	Sender   string
	Receiver string
	Date     string
	TS       TimeStamp
}

func NewRequest(sender string, date string, timeStamp TimeStamp) *Message {
	return &Message{
		MsgType: Request,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		Date: date,
		TS:   timeStamp,
	}
}

func NewReply(sender string, receiver string, date string, timeStamp TimeStamp) *Message {
	return &Message{
		MsgType:  Reply,
		Sender:   sender,
		Receiver: receiver,
		Date:     date,
		TS:       timeStamp,
	}
}

func NewRelease(sender string, date string, timeStamp TimeStamp) *Message {
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

func WriteMsgToFile(path string, id string, typeMsg string, message Message, timestamp TimeStamp) error {

	var err error
	f := utilities.OpenFile(path)
	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)
	if typeMsg == "send" {
		_, err = f.WriteString("[" + date + "] : " + id + " sends" + message.ToString("send") + " to " + message.Receiver + ".")
	}
	if typeMsg == "receive" {
		_, err = f.WriteString("[" + date + "] : " + id + " receives" + message.ToString("receive"))
		if message.MsgType != Reply {
			_, err = f.WriteString(" and updates its local logical timestamp to " + strconv.Itoa(int(timestamp)))
		}
	}
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
	return nil
}
