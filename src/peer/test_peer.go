package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
)

func runTest() {

	setAlgorithmPeer()

	switch algorithm {
	case "tokenAsking":

		tokenAsking.ExecuteTestPeer(&myTokenPeer, numSenders)

	case "ricartAgrawala":
		ricartAgrawala.ExecuteTestPeer(&myRApeer, numSenders)
	case "lamport":
		lamport.ExecuteTestPeer(&myLamportPeer, numSenders)

	}

}
