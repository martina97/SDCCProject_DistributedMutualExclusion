package utilities

import (
	"log"
	"os"
	"time"
)

func WriteInfosToFile(text string, path string, username string) {

	f := OpenFile(path)

	//save new address on file
	date := time.Now().Format(DateFormat)

	_, _ = f.WriteString("[" + date + "] : " + username + " " + text)

	_, _ = f.WriteString("\n")
	_ = f.Sync()
}

func OpenFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}

func WriteVCInfoToFile(path string, username string, vcString string) {

	f := OpenFile(path)

	//save new address on file
	date := time.Now().Format(DateFormat)

	_, err := f.WriteString("[" + date + "] : " + username + " update its vector clock to " + vcString)
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}

func WriteTSInfoToFile(path string, id string, timestamp string) {

	f := OpenFile(path)

	//save new address on file
	date := time.Now().Format(DateFormat)

	_, err := f.WriteString("[" + date + "] : " + id + " updates its local logical timeStamp to " + timestamp)
	_, err = f.WriteString("\n")
	err = f.Sync()
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}
