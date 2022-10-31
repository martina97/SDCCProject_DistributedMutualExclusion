package utilities

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"strconv"
)

/*
	Set info
*/
func setInfo(info *Process, port int, username string) error {
	info.Address = GetLocalIP()
	if info.Address == "" {
		return errors.New("Impossible to find local ip")
	}

	info.Port = strconv.Itoa(port)

	info.Username = username
	return nil
}

/*
	Registration function for peer
*/
func Registration(peers *list.List, port int, username string, listNodes []Process) Result_file {

	var info Process
	var res Result_file

	fmt.Println("SONO IN REGISTRATION --- PORTA == ", port)

	addr := Server_addr + ":" + strconv.Itoa(Server_port)
	fmt.Println("SONO IN REGISTRATION --- ADDR	 == ", addr)

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
	log.Printf("Call to registration node")
	err = server.Call("Utility.Save_registration", &info, &res)
	if err != nil {
		log.Fatal("Error save_registration procedure: ", err)
	}

	//check result
	for e := 0; e < res.PeerNum; e++ {
		var item Process
		var tmp string
		item.Username, tmp, item.Address, item.Port = ParseLine(res.Peers[e], ":")
		fmt.Println("tmp ======= ", tmp)
		fmt.Println("item.Address ======= ", item.Address)
		fmt.Println(" item.Port ======= ", item.Port)

		fmt.Println("res.Peers[e] ======= ", res.Peers[e])
		item.Type = StringToType(tmp)
		item.ID = e //setto id del peer

		peers.PushBack(item)
	}
	/*
		res contiene Result_file con tutte le info sul file di log che sono state memorizzate da Save_registration
		perche gli ho passato come parametro Result_file
	*/
	fmt.Println("PROVA STAMPA RES	 == ", res)

	return res
}
