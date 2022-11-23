package utilities

const (
	MAXCONNECTION int    = 3 //MAXPEERS + 1 sequencer
	COORDINATOR   string = "p0"
	MAXPEERS      int    = 3
	Server_port   int    = 4321
	Server_addr   string = "10.10.1.50"
	// Server_addr string = "localhost" //if running outside docker
	Client_port       int    = 2345
	Peer_msg_sca_file string = "/docker/node_volume/messageSca.txt"
	Clean_Test_Dir    bool   = true
	MSG_BUFFERED_SIZE int    = 100
	CHAN_SIZE         int    = 1
	DATE_FORMAT       string = "15:04:05.000"
	Test              bool   = true
)
