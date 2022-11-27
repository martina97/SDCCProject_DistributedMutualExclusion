package utilities

import (
	"log"
	"os"
	"strings"
)

type NodeType int

const (
	Peer     NodeType = 0
	Register          = 1
)

// NodeInfo : infos about node
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

func CreateLog(path string, header string) *log.Logger {
	serverLogFile, err := os.Create(path)
	if err != nil {
		log.Printf("unable to read file: %v", err)
	}
	serverLogFile.Close()
	return log.New(serverLogFile, header, log.Lshortfile)
}
