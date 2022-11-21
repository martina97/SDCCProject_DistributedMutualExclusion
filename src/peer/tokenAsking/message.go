package tokenAsking

import "SDCCProject_DistributedMutualExclusion/src/utilities"

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
