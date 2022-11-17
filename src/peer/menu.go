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
			Items: []string{"Run centralized", "Run Lamport", "Run Ricart-Agrawala", "exit"},
		}

		i, result, err := prompt.Run() //i : indice di cosa ho scelto

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose number %d: %s\n", i+1, result)

		if i+1 == 1 { //ossia Run centralized
			//utilities.Registration(peers, utilities.Client_port, username, listNodes)
			algorithm = "centralized"
			openSecondMenu()
			setAlgorithmPeer()
			//openCentralizedMenu()
		}
		if i+1 == 2 { //ossia Run Lamport
			//utilities.Registration(peers, utilities.Client_port, username, listNodes)
			algorithm = "Lamport"
			openSecondMenu()
			//openLamportMenu()
		}
		if i+1 == 3 { //exit
			algorithm = "RicartAgrawala"
			setAlgorithmPeer()
			setPeerUtils2()
			openSecondMenu()

		}
		if i+1 == 4 { //exit
			break
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
			Items: []string{"send message", "show message received", "exit"},
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

func openCentralizedMenu() {

	for { //infinite loop
		prompt := promptui.Select{
			Label: "Select Option",
			Items: []string{"send message", "show message received", "exit"},
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

func openLamportMenu() {

	//creo ora la lista di msg del peer

	for { //infinite loop
		prompt := promptui.Select{
			Label: "Select Option",
			Items: []string{"send message", "show message received", "exit"},
		}

		i, result, err := prompt.Run() //i : indice di cosa ho scelto

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose number %d: %s\n", i+1, result)

		if i+1 == 1 { //Lamport
			sendMessage()
		}

		if i+1 == 3 { //exit
			break
		}

	}
}

//menu per decidere se, quando ricevo msg REQUEST, mandare REPLY o NEW REQUEST, ossia se inviare ack oppure nuova richiesta
func openMenuRequest() string {

	//for { //infinite loop
	prompt := promptui.Select{
		Label: "Select Option",
		Items: []string{"send ACK", "send NEW REQUEST"},
	}

	i, result, err := prompt.Run() //i: indice di cosa ho scelto

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	fmt.Printf("You choose number %d: %s\n", i+1, result)

	if i+1 == 1 { //send ack
		return "ack"
	}

	if i+1 == 2 { //send new request
		return "ack"
	}
	return ""

	//}
}
