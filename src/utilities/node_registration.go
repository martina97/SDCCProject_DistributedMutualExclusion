package utilities

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"strconv"
)

func setInfo(info *NodeInfo, port int, username string) error {
	info.Address = GetLocalIP()
	if info.Address == "" {
		return errors.New("Impossible to find local ip")
	}

	info.Port = strconv.Itoa(port)

	info.Username = username
	return nil
}

func Registration(peers *list.List, port int, username string) Result_file {

	var info NodeInfo
	var res Result_file

	addr := ServerAddr + ":" + strconv.Itoa(ServerPort)

	// Try to connect to addr
	server, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer server.Close()

	//set info to send
	err = setInfo(&info, port, username)
	if err != nil {
		log.Fatal("Error on setInfo: ", err)
	}

	//call procedure
	fmt.Println("Call to register")
	err = server.Call("Utility.SaveRegistration", &info, &res)
	if err != nil {
		log.Fatal("Error SaveRegistration procedure: ", err)
	}

	//check result
	for e := 0; e < res.PeerNum; e++ {
		var item NodeInfo
		var tmp string
		item.Username, tmp, item.Address, item.Port = ParseLine(res.Peers[e], ":")

		item.Type = StringToType(tmp)
		item.ID = e //setto id del peer

		peers.PushBack(item)
	}
	/*
		res contiene Result_file con tutte le info sul file di log che sono state memorizzate da Save_registration
		perche gli ho passato come parametro Result_file
	*/

	return res
}
