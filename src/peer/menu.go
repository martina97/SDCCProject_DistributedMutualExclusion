package main

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala"
	"SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking"
	"fmt"
	"github.com/manifoldco/promptui"
)

func openMenu() {

	for { //infinite loop
		prompt := promptui.Select{
			Label: "Select Option",
			Items: []string{"Run Token-Asking", "Run Lamport", "Run Ricart-Agrawala"},
		}

		i, result, err := prompt.Run() //i: indice di cosa ho scelto

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose number %d: %s\n", i+1, result)

		if i+1 == 1 { //tokenAsking
			algorithm = "tokenAsking"
			setAlgorithmPeer()
			openSecondMenu()

		}
		if i+1 == 2 { //Lamport
			algorithm = "lamport"
			setAlgorithmPeer()
			openSecondMenu()
		}
		if i+1 == 3 { //ricartAgrawala
			algorithm = "ricartAgrawala"
			setAlgorithmPeer()
			openSecondMenu()
		}
	}
}

func openSecondMenu() {

	for { //infinite loop
		prompt := promptui.Select{
			Label: "Select Option",
			Items: []string{"send request", "exit"},
		}

		i, result, err := prompt.Run() //i: indice di cosa ho scelto

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose number %d: %s\n", i+1, result)

		if i+1 == 1 { //send message
			sendMessageRequest()
		}

		if i+1 == 2 { //exit
			break
		}
	}

}

func openTestMenu() (int, string) {
	prompt := promptui.Select{
		Label: "Select Option",
		Items: []string{"Test TokenAsking 1 sender", "Test TokenAsking 2 senders",
			"Test RicartAgrawala 1 sender", "Test RicartAgrawala 2 senders",
			"Test Lamport 1 sender", "Test Lamport 2 senders"},
	}

	i, result, err := prompt.Run() //i: indice di cosa ho scelto

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	fmt.Printf("You choose number %d: %s\n", i+1, result)

	if i+1 == 1 { //send ack
		return 1, "tokenAsking"
	}

	if i+1 == 2 { //send new request
		return 2, "tokenAsking"
	}
	if i+1 == 3 { //send new request
		return 1, "ricartAgrawala"
	}
	if i+1 == 4 { //send new request
		return 2, "ricartAgrawala"
	}
	if i+1 == 5 { //send new request
		return 1, "lamport"
	}
	if i+1 == 6 { //send new request
		return 2, "lamport"
	}
	return 0, ""
}

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
