package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

func startTests() {

	runTest(1, "tokenAsking") //invia msg solo il peer con ID 1
	//runTest(2, "tokenAsking") //invia msg solo il peer con ID 1
	runTest(1, "ricartAgrawala") //invia msg solo il peer con ID 1

}

func runTest(i int, s string) {
	fmt.Println("sto in runTest")
	fmt.Println(myTokenPeer)
	fmt.Println("algorithm == ", algorithm)

	algorithm = s
	setAlgorithmPeer()
	fmt.Println(myTokenPeer)

	switch s {
	case "tokenAsking":
		if myUsername != utilities.COORDINATOR {
			tokenAsking.ExecuteTestPeer(&myTokenPeer, i)
		} else {
			tokenAsking.ExecuteTestCoordinator(&myCoordinator, i)
		}
	case "ricartAgrawala":
		ricartAgrawala.ExecuteTestPeer(&myRApeer, i)

	}

}
