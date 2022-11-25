package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

func startTests() {

	runTest()
	//runTest(1, "tokenAsking") //invia msg solo il peer con ID 1
	//runTest(2, "tokenAsking") //invia msg solo il peer con ID 1
	//runTest(1, "ricartAgrawala") //invia msg solo il peer con ID 1

}

//func runTest(i int, s string) {
func runTest() {
	fmt.Println("sto in runTest: ", algorithm, "num sender = ", numSenders)
	fmt.Println(myTokenPeer)
	fmt.Println("algorithm == ", algorithm)

	//algorithm = s
	setAlgorithmPeer()
	fmt.Println(myTokenPeer)

	switch algorithm {
	case "tokenAsking":
		if myUsername != utilities.COORDINATOR {
			tokenAsking.ExecuteTestPeer(&myTokenPeer, numSenders)
		} else {
			tokenAsking.ExecuteTestCoordinator(&myCoordinator, numSenders)
		}
	case "ricartAgrawala":
		ricartAgrawala.ExecuteTestPeer(&myRApeer, numSenders)

	}

}
