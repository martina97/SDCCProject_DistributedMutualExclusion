package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
	"github.com/manifoldco/promptui"
)

func openMenu() {

	myNode.ScalarMap = utilities.MessageMap{} //inizializzo mappa
	for {                                     //infinite loop
		prompt := promptui.Select{
			Label: "Select Option",
			Items: []string{"Run Token-Asking", "Run Lamport", "Run Ricart-Agrawala"},
		}

		i, result, err := prompt.Run() //i : indice di cosa ho scelto

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose number %d: %s\n", i+1, result)

		if i+1 == 1 { //ossia Run tokenAsking
			//utilities.Registration(peers, utilities.Client_port, username, listNodes)
			algorithm = "tokenAsking"

			//qui, in base al fatto se sono coordinatore o meno, ho 2 diverse cose, infatti va in openSecondMenu solo
			//chi non è coordinatore
			//decido che il coordinatore è p0
			if myUsername == utilities.COORDINATOR {
				fmt.Println("sono il coordinatore")
			} else {
				fmt.Println("non sono il coordinatore")
			}
			setAlgorithmPeer()
			if myUsername != utilities.COORDINATOR {
				openSecondMenu()
			}
		}
		if i+1 == 2 { //ossia Run Lamport
			//utilities.Registration(peers, utilities.Client_port, username, listNodes)
			algorithm = "lamport"
			setAlgorithmPeer()

			openSecondMenu()
			//openLamportMenu()
		}
		if i+1 == 3 { //ricartAgrawala
			algorithm = "ricartAgrawala"
			setAlgorithmPeer()
			//setPeerUtils2()
			openSecondMenu()

		}

	}

}

func openSecondMenu() {

	// una volta scelto l'algoritmo, setto le info dei vari peer (in particolare il file di log, il cui path
	// dipende dall0'algoritmo scelto

	fmt.Println(" sto in openSecondMenu ------ ")
	for { //infinite loop
		prompt := promptui.Select{
			Label: "Select Option",
			Items: []string{"send request", "show message received", "exit"},
		}

		i, result, err := prompt.Run() //i : indice di cosa ho scelto

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose number %d: %s\n", i+1, result)

		if i+1 == 1 { //send message
			sendMessage()
		}

		if i+1 == 3 { //exit
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
