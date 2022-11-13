package utilities

const (
//MSG_BUFFERED_SIZE = 100
//CHAN_SIZE         = 1
)

/*
// Struct to send information about lock
// mi serve per tenere info sui msg ricevuti, lock, lista msg ...
//todo: rinominare con process
type InfoLock struct {
	//TODO: servono? o basta il processo che ha gia tali info?
	//nodeID   int
	//nodePort string

	//algorithm
	// algorithim
	Waiting         bool //serve a vedere se chi ha mandato msg request e' in attesa di tutti i msg reply
	chanRcvMsg      chan Message
	chanSendMsg     chan *Message
	ChanAcquireLock chan bool
	ReplyProSet     *list.List // then Message.Sender is the key.
	deferProSet     *list.List // then Message.Sender is the key.
	//msgMap          *MessageMap // request lock message priority queue.

	// process handler
	//proc *peer.process
	peer *NodeInfo
	// log
	logger *log.Logger

	mu sync.Mutex

	//peer
	//tcpPeer NodeInfo

	//timestamp

}

*/

/*
func (l *InfoLock) GetMutex() *sync.Mutex {
	// p.mu.Lock()
	// defer p.mu.Unlock()
	return &l.mu
}

*/
/*
func NewLock(peer *NodeInfo) (*InfoLock, error) {
	fmt.Println("########   SONO IN NEW LOCK 	##################### \n\n ")
	dl := &InfoLock{
		//todo: metto peer id e nodeport o basta info proc?
		chanRcvMsg:      make(chan Message, MSG_BUFFERED_SIZE),
		chanSendMsg:     make(chan *Message, MSG_BUFFERED_SIZE),
		ChanAcquireLock: make(chan bool, CHAN_SIZE),
		ReplyProSet:     list.New(),
		deferProSet:     list.New(),
		//proc:            process,
		peer: peer,
		//msgMap: &MessageMap{},
	}
	//fmt.Println("dl.chanAcquireLock ==", reflect.TypeOf(dl.ChanAcquireLock))
	dl.logger = CreateLog("lockInfo_", strconv.Itoa(peer.ID), "[infoLock] ")

	//metto peer in ascolto
	//utilities.StartListen(process.tcpPeer)

	f, err := os.OpenFile("/docker/node_volume/lockInfo_"+strconv.Itoa(peer.ID)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	dl.logger.SetOutput(f)
	dl.logger.Println("infoLock(" + strconv.Itoa(dl.peer.ID) + ") created.\n")
	//dl.tcpPeer = peer
	//peer.LockInfo = dl

	//qui devo fare handle connection
	fmt.Println("dl ===", dl)
	return dl, nil
}

*/

/*
func CreateLog(typeInfo string, id string, header string) *log.Logger {
	serverLogFile, err := os.Create("/docker/node_volume/" + typeInfo + id + ".log")
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

//return log.New(serverLogFile, header, log.Lshortfile)
//}
