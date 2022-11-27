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
