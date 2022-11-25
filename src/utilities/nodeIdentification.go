package utilities

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NodeType int

const (
	Peer     NodeType = 0
	Register          = 1
)
const (
	//MSG_BUFFERED_SIZE  = 100
	//CHAN_SIZE          = 1
	CONN_BUFFERED_SIZE = 1
)

// Struct to send information about peer
// processo in esecuzione sul peer
type NodeInfo struct {

	//info su nodo su cui e' in esecuzione il processo
	Username string   //nome nodo
	Type     NodeType // tipo di nodo
	ID       int      //id nodo
	Address  string   //indirizzo nodo
	Port     string   //porta nodo

	// utili per mutua esclusione
	mutex     sync.Mutex
	timestamp TimeStamp
	fileLog   *log.Logger //file di ogni processo in cui scrivo info di quando accede alla sez critica
	Listener  net.Listener

	// algorithim
	Waiting         bool //serve a vedere se chi ha mandato msg request e' in attesa di tutti i msg reply
	ChanRcvMsg      chan Message
	ChanSendMsg     chan *Message
	ChanAcquireLock chan bool
	replyProSet     *list.List // then Message.Sender is the key.
	deferProSet     *list.List // then Message.Sender is the key.

	ScalarMap MessageMap
	TimeStamp TimeStamp
	LogPath   string
	//LockInfo *infoLock
}

func TypeToString(nodeType NodeType) string {
	switch nodeType {
	case Peer:
		return "peer"
	case Register:
		return "register"
	}
	return ""
}

func StringToType(s string) NodeType {
	switch s {
	case "peer":
		return Peer
	case "register":
		return Register
	}
	return -1
}

func ParseLine(s string, sep string) (string, string, string, string) {
	res := strings.Split(s, sep)
	return res[0], res[1], res[2], res[3]
}

func CreateLog2(path string, header string) *log.Logger {
	serverLogFile, err := os.Create(path)
	if err != nil {
		log.Printf("unable to read file: %v", err)
	}
	serverLogFile.Close()
	/*
		newpath := filepath.Join(".", "log")
		os.MkdirAll(newpath, os.ModePerm)
		serverLogFile, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		// return log.New(serverLogFile, header, log.Lmicroseconds|log.Lshortfile)

	*/

	//process.SetFileLog(log.New(serverLogFile, header, log.Lshortfile))

	return log.New(serverLogFile, header, log.Lshortfile)
}

func WriteMsgToFile3(path string, id string, typeMsg string, message Message, timestamp TimeStamp, algo string) error {
	fmt.Println("sto in WriteMsgToFile3")
	fmt.Println("path == ", path)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//save new address on file
	date := time.Now().Format(DATE_FORMAT)
	if typeMsg == "send" {
		_, err = f.WriteString("[" + date + "] : " + id + " " + typeMsg + message.ToString("send") + " to " + message.Receiver + ".")
	}
	if typeMsg == "receive" {
		_, err = f.WriteString("[" + date + "] : " + id + " " + typeMsg + message.ToString("receive"))
		if message.MsgType != Reply { //in ricart il TS lo aggiorno solo quando ricevo REQUEST
			_, err = f.WriteString(" and update its local logical timestamp to " + strconv.Itoa(int(timestamp)))
		}
	}
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
	return nil
}

func WriteInfoToFile2(username string, path string, text string, infoCS bool) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//save new address on file
	date := time.Now().Format(DATE_FORMAT)

	if infoCS == false {
		_, err = f.WriteString("[" + date + "] : " + username + " " + text)
	} else {
		_, err = f.WriteString("\n" + username + text)

	}
	_, err = f.WriteString("\n")
	err = f.Sync()
}

func WriteTSInfoToFile2(path string, id string, timestamp TimeStamp, algorithm string) {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	//save new address on file
	date := time.Now().Format(DATE_FORMAT)

	_, err = f.WriteString("[" + date + "] : " + id + " update its local logical timeStamp to " + strconv.Itoa(int(timestamp)))
	_, err = f.WriteString("\n")
	err = f.Sync()
}

func WriteVCInfoToFile(path string, id string, vc VectorClock) {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	//save new address on file
	date := time.Now().Format(DATE_FORMAT)

	_, err = f.WriteString("[" + date + "] : " + id + " update its local vector clock to " + ToString(vc))
	_, err = f.WriteString("\n")
	err = f.Sync()
}
