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
/*
type CentralizedMessage struct {
	//MsgID   string
	MsgTypeCentr MessageTypeCentr //request,reply,release
	Sender       Process
	Receiver     int
	Date         string
}

//// ############################ ALGORITMO CENTRALIZZATO #######################################

func EnterMsg(process Process, date string) *CentralizedMessage {
	return &CentralizedMessage{
		MsgTypeCentr: Enter,
		Sender:       process,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		//MsgContent: msgContent,
		Date: date,
	}
}

func GrantedMsg(receiver int, date string) *CentralizedMessage {
	return &CentralizedMessage{
		MsgTypeCentr: Granted,
		Receiver:     receiver,
		//Receiver:   receiver,	//non serve specificarlo perche la richiesta viene mandata a tutti
		//MsgContent: msgContent,
		Date: date,
	}
}

*/
