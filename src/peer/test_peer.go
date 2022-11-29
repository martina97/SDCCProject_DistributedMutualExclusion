package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"SDCCProject_DistributedMutualExclusion/src/utilities"
)

func runTest() {

	setAlgorithmPeer()

	switch algorithm {
	case "tokenAsking":
		if myNode.Username != utilities.COORDINATOR {
			tokenAsking.ExecuteTestPeer(&myTokenPeer, numSenders)
		} else {
			//tokenAsking.ExecuteTestCoordinator(&myCoordinator, numSenders)
		}
	case "ricartAgrawala":
		ricartAgrawala.ExecuteTestPeer(&myRApeer, numSenders)
	case "lamport":
		lamport.ExecuteTestPeer(&myLamportPeer, numSenders)

	}

}
