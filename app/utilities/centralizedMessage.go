package utilities

const (
	Enter   MessageType = iota + 1 // request mutual lock
	Granted                        // reply mutual lock
	Denied                         // release mutual lock ///todo: lo levo xk lo faccio sincrono

)

/*
servono lettera maiuscola perche devo farne il marshaling e unmarshaling, che richiedono solo field esportabili
todo: vedi libro pag 108
*/
type CentralizedMessage struct {
	//MsgID   string
	MsgType  MessageType //request,reply,release
	Sender   int
	Receiver int
	Date     string
}

//// ############################ ALGORITMO CENTRALIZZATO #######################################

func EnterMsg(sender int, date string) *Message {
	return &Message{
		//MsgID:      RandStringBytes(msgIDCnt),
		MsgType: Request,
		Sender:  sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		//MsgContent: msgContent,
		Date: date,
	}
}
