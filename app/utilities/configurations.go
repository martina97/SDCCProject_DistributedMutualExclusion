package utilities

const (
	MAXCONNECTION int    = 3 //MAXPEERS + 1 sequencer
	MAXPEERS      int    = 3
	Server_port   int    = 4321
	Server_addr   string = "10.10.1.50"
	// Server_addr string = "localhost" //if running outside docker
	Client_port       int    = 2345
	Peer_msg_sca_file string = "/docker/node_volume/messageSca.txt"
	Launch_Test       bool   = false //launch all peer in test mode
	Clean_Test_Dir    bool   = true
	MSG_BUFFERED_SIZE int    = 100
	CHAN_SIZE         int    = 1
)
