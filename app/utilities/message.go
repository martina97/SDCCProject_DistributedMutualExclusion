package utilities

import (
	"container/list"
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
	Sender   int
	Receiver int
	SeqNum   []uint64
	Date     string
	TS       TimeStamp
}

// NewRequest returns a new distributed mutual lock message.
func NewRequest(ts []uint64, sender int, date string, timeStamp TimeStamp) *Message {
	return &Message{
		//MsgID:      RandStringBytes(msgIDCnt),
		MsgType: Request,
		SeqNum:  ts,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		//MsgContent: msgContent,
		Date: date,
		TS:   timeStamp,
	}
}

// NewRequest returns a new distributed mutual lock message.
func NewRequest2(sender int, date string, timeStamp TimeStamp) *Message {
	return &Message{
		//MsgID:      RandStringBytes(msgIDCnt),
		MsgType: Request,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		//MsgContent: msgContent,
		Date: date,
		TS:   timeStamp,
	}
}

func NewReply(ts []uint64, sender int, receiver int, date string, timeStamp TimeStamp) *Message {
	return &Message{
		SeqNum:   ts,
		MsgType:  Reply,
		Sender:   sender,
		Receiver: receiver,
		Date:     date,
		TS:       timeStamp,
	}
}
func NewRelease(sender int, date string, timeStamp TimeStamp) *Message {
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

func (m *Message) MessageToString(role string) string {
	var name string
	//date := time.Now().Format("2006/01/02 15:04:05")

	switch m.MsgType {
	case Request:
		name = "Request"
	case Release:
		name = "Release"
	case Reply:
		name = "Reply"
	}

	if role == "send" {
		return fmt.Sprintf(" %s message: {%s %d %d %d %s [%d]}", name, name, m.SeqNum, m.Sender, m.Receiver, m.Date, m.TS)
	}
	if role == "receive" {
		return fmt.Sprintf(" %s message: {%s %d %d %d %s [%d]} from process(%d)", name, name, m.SeqNum, m.Sender, m.Receiver, m.Date, m.TS, m.Sender)
	}

	return ""
}

func InsertInOrder(l *list.List, msg Message) *list.List {
	tmp := msg.SeqNum[0]

	fmt.Println("SONO IN InsertInOrder, MSG === ", msg)
	//scorro lista msg gia presenti
	for e := l.Front(); e != nil; e = e.Next() {
		item := e.Value.(Message)

		if tmp < item.SeqNum[0] { //msg ha TS minore del primo in lista
			l.InsertBefore(msg, e)
			return l
		}
	}
	l.PushBack(msg)
	return l
}
