package utilities

import (
	"bufio"
	"container/list"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

/*
	This package static configure connection port and ip
*/

type Utility int

// Constant value

/*
	Non optimal solution:
	MAXCONNECTION = numberOfPeer + 1 (sequencer)
*/

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)

	// per algo centralizzato
	resourceState bool = true // è true: risorsa libera, false: risorsa occupata
	queue              = list.New()
)

type Result_file struct {
	PeerNum int
	Peers   [MAXCONNECTION]string
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// save registration info to reg_node procedure
func (utils *Utility) Save_registration(arg *Process, res *Result_file) error {

	log.Printf("The registration is for %s the ip address:port : %s:%s\n", TypeToString(arg.Type), arg.Address, arg.Port)
	fmt.Println("provaaaaaa")
	/*
		myfile, e := os.Create("provaFILEEEEE.txt")
		if e != nil {
			log.Fatal(e)
		}
		log.Println(myfile)
		myfile.Close()
	*/
	f, err := os.OpenFile("/docker/register_volume/nodes.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Println(err)
		return errors.New("Impossible to open file")
	}
	fmt.Println("provaaaaaa2")

	/*
		see https://www.joeshaw.org/dont-defer-close-on-writable-files/ for file defer on close
	*/
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	//save new address on file
	_, err = f.WriteString(arg.Username + ":" + TypeToString(arg.Type) + ":" + arg.Address + ":" + arg.Port)
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		return err
	}

	log.Printf("Saved")

	Connection <- true
	Wg.Add(1)
	log.Printf("Waiting other connection")
	Wg.Wait()

	//send back file
	err = prepare_response(res)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (utils *Utility) CentralizedSincro(arg *Message, res *Result_file) error {
	log.Printf("sono in CentralizedSincro")
	log.Printf("*arg == ", *arg)
	log.Printf("&arg == ", &arg)

	if arg.MsgType == Request { //msg di request
		if resourceState == true {
			// processo puo accedere in CS
			log.Printf("processo puo accedere in CS")

			log.Printf(" arg.SenderProc.Address == ", &arg.SenderProc.Address)
			log.Printf(" arg.SenderProc.Port == ", &arg.SenderProc.Port)

			//devo inviare msg granted al processo
			peerConn := arg.SenderProc.Address + ":" + arg.SenderProc.Port

			conn, err := net.Dial("tcp", peerConn)
			defer conn.Close()
			if err != nil {
				log.Println("Send response error on Dial")
			}
			//msg2 := new(CentralizedMessage)

			/*
				dec := gob.NewDecoder(conn)
				dec.Decode(msg2)
				dec.Decode(arg)
				log.Printf("sono DENTRO IF")
				log.Printf("*msg2 == ", *msg2)
				log.Printf("&msg2 == ", &msg2)
				log.Printf("*arg == ", *arg)
				log.Printf("&arg == ", &arg)

			*/

			date := time.Now().Format("15:04:05.000")

			msg := NewReply2(3, arg.Sender, date, 0)
			enc := gob.NewEncoder(conn)
			enc.Encode(msg)

			resourceState = false
		} else {

		}
	} else if arg.MsgType == Release {
		log.Printf("msg di release")

	}

	return nil
}

func prepare_response(res *Result_file) error {
	res.PeerNum = MAXCONNECTION
	file, err := os.Open("/docker/register_volume/nodes.txt")
	if err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var i int
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) //stampa contenuto file
		// manda ai peer chiamando peer.saveFile

		res.Peers[i] = line
		i++
	}
	if err := scanner.Err(); err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	err = file.Sync()
	if err != nil {
		return errors.New("error on open file[prepare_file]")
	}
	return nil
}

// This function allow to verify if there is error and return it.
func Check_error(err error) error {
	if err != nil {
		log.Printf("unable to read file: %v", err)
	}
	return err
}
