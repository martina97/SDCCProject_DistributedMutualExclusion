package utilities

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type NodeType int

const (
	Peer     NodeType = 0
	Register          = 1
)

// NodeInfo : infos about peer
type NodeInfo struct {

	//info su nodo su cui e' in esecuzione il processo
	Username string   //peer name
	Type     NodeType //peer type
	ID       int      //peer ID
	Address  string   //node address
	Port     string   //node port

	LogPath string
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
