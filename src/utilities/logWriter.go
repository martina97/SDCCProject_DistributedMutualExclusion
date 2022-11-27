package utilities

import (
	"log"
	"os"
	"time"
)

func WriteInfosToFile(text string, path string, username string) {

	f := openFile(path)

	//save new address on file
	date := time.Now().Format(DATE_FORMAT)

	_, _ = f.WriteString("[" + date + "] : " + username + " " + text)

	_, _ = f.WriteString("\n")
	_ = f.Sync()
}

func openFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}

func WriteVCInfoToFile(path string, isCoord bool) {
	var vc VectorClock
	var username string

	f := openFile(path)

	if isCoord {
		vc = myCoordinator.VC
		username = "coordinator"
	} else {
		vc = myPeer.VC
		username = myPeer.Username
	}

	//save new address on file
	date := time.Now().Format(utilities.DATE_FORMAT)

	_, err := f.WriteString("[" + date + "] : " + username + " update its vector clock to " + ToString(vc))
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}
