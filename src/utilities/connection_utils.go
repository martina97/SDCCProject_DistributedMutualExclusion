package utilities

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

/*
	This package static configure connection port and ip
*/

type Utility int

var (
	Connection = make(chan bool)
	Wg         = new(sync.WaitGroup)
)

type Result_file struct {
	PeerNum int
	Peers   [MAXPEERS]string
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

// SaveRegistration : save registration infos
func (utils *Utility) SaveRegistration(arg *NodeInfo, res *Result_file) error {

	log.Printf("The registration is for %s the ip address:port : %s:%s\n", TypeToString(arg.Type), arg.Address, arg.Port)

	f, err := os.OpenFile("/docker/register_volume/nodes.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Println(err)
		return errors.New("impossible to open file")
	}

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
	err = prepareResponse(res)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func prepareResponse(res *Result_file) error {
	res.PeerNum = MAXPEERS
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

// CheckError : This function allow to verify if there is error and return it.
func CheckError(err error, text string) {
	if err != nil {
		log.Fatalf("%s: %v", text, err)
	}
}

func SleepRandInt() {
	min := 300
	max := 2000
	// set seed
	rand.Seed(time.Now().UnixNano())
	// generate random number and print on console
	delay := rand.Intn(max-min) + min
	//fmt.Println(delay)
	time.Sleep(time.Millisecond * time.Duration(delay))
}
