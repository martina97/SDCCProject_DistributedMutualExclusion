package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

func startTests() {

	/*
		runTest(1, "tokenAsking") //invia msg solo il peer con ID 1
		fmt.Println("runTest1 finito")

	*/
	runTest(2, "tokenAsking") //invia msg solo il peer con ID 1

}

func runTest(i int, s string) {
	fmt.Println("sto in runTest")
	fmt.Println(myTokenPeer)
	fmt.Println("algorithm == ", algorithm)

	algorithm = "tokenAsking"
	setAlgorithmPeer()
	fmt.Println(myTokenPeer)

	if myUsername != utilities.COORDINATOR {
		tokenAsking.ExecuteTestPeer(&myTokenPeer)
	} else {
		tokenAsking.ExecuteTestCoordinator(&myCoordinator)
	}

	/*
		if i == 1 {
			if myTokenPeer.Username == "p2" {
				err := sendMessage()
				if err != nil {
					log.Fatalf("error sending request %v", err)
				}
			} else {
				//time.Sleep(time.Minute / 2)
				//select {}
				time.Sleep(time.Minute)
			}
		} else {
			if myTokenPeer.Username == "p1" || myTokenPeer.Username == "p2" {
				err := sendMessage()
				if err != nil {
					log.Fatalf("error sending request %v", err)
				}
			} else {
				//time.Sleep(time.Minute / 2)
				//select {}
				time.Sleep(time.Minute + time.Minute/2)
			}
		}

	*/

}
