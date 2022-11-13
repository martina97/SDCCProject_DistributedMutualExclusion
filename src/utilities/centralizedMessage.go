package utilities

import "fmt"

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
	Sender       NodeInfo
	Receiver     int
	Date         string
}

//// ############################ ALGORITMO CENTRALIZZATO #######################################

func NewEnterMsg(process NodeInfo, date string) *CentralizedMessage {
	return &CentralizedMessage{
		MsgTypeCentr: Enter,
		Sender:       process,
		Date:         date,
	}
}

func NewGrantedMsg(receiver int, date string) *CentralizedMessage {
	return &CentralizedMessage{
		MsgTypeCentr: Granted,
		Receiver:     receiver,
		Date:         date,
	}
}

func (m *CentralizedMessage) MessageToString(role string) string {
	var name string
	//date := time.Now().Format("2006/01/02 15:04:05")

	switch m.MsgTypeCentr {
	case Enter:
		name = "Request"
	case Granted:
		name = "Granted"
	case Denied:
		name = "Denied"
	}

	if role == "send" {
		return fmt.Sprintf(" %s message: {%s %d %d %s}", name, name, m.Sender.ID, m.Receiver, m.Date)
		return fmt.Sprintf(" %s message: {%s %d %s}", name, m.Sender.ID, m.Receiver, m.Date)
	}
	if role == "receive" {
		return fmt.Sprintf(" %s message: {%s %d %d %s} from %d", name, name, m.Sender.ID, m.Receiver, m.Date, m.Sender.ID)
	}

	return ""
}
