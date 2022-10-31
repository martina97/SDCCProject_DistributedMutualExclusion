package utilities

type MessageTypeCentr string

const (
	Enter   MessageTypeCentr = "enter"   // request mutual lock
	Granted MessageTypeCentr = "granted" // reply mutual lock
	Denied  MessageTypeCentr = "denied"  // release mutual lock ///todo: lo levo xk lo faccio sincrono

)

/*
servono lettera maiuscola perche devo farne il marshaling e unmarshaling, che richiedono solo field esportabili
todo: vedi libro pag 108
*/
type CentralizedMessage struct {
	//MsgID   string
	MsgTypeCentr MessageTypeCentr //request,reply,release
	Sender       int
	Receiver     int
	Date         string
}

//// ############################ ALGORITMO CENTRALIZZATO #######################################

func EnterMsg(sender int, date string) *CentralizedMessage {
	return &CentralizedMessage{
		MsgTypeCentr: Enter,
		Sender:       sender,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		//MsgContent: msgContent,
		Date: date,
	}
}
