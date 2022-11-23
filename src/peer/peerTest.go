package main

import (
	"fmt"
	"log"
)

func startTests() {

	runTest(1, "TokenAsking") //invia msg solo il peer con ID 1
}

func runTest(i int, s string) {
	fmt.Println("sto in runTest")
	fmt.Println(myTokenPeer)
	fmt.Println("algorithm == ", algorithm)

	algorithm = "tokenAsking"
	setAlgorithmPeer()
	fmt.Println(myTokenPeer)
	if myTokenPeer.Username == "p1" {
		err := sendMessage()
		if err != nil {
			log.Fatalf("error sending request %v", err)
		}
	} else {
		//time.Sleep(time.Minute / 2)
		select {}
	}

}
