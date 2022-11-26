package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
)

//send msg di request
func sendMessageRequest() {
	switch algorithm {
	case "tokenAsking":
		tokenAsking.SendRequest(&myTokenPeer)
	case "lamport":
		lamport.SendLamport(&myLamportPeer)
	case "ricartAgrawala":
		ricartAgrawala.SendRicart(&myRApeer)
	}
}
