package utilities

import (
	"SDCCProject_DistributedMutualExclusion/src/peer/lamport"
	"fmt"
)

// MessageMap : lista di messaggi, ogni messaggio ha timestamp, id, tipo, receiver e sender todo: cambiare
type MessageMap map[TimeStamp][]lamport.Message

// test inserts several ints into an MessageHeap, checks the minimum,
// and removes them in order of priority.

func AppendHashMap(map1 MessageMap, message lamport.Message) {
	var listMsg []lamport.Message
	//lista := list.New()
	//lista.PushBack(message)

	//fmt.Println("map1 ===", map1)
	fmt.Println("msg ==", message)
	_, ok := map1[message.TS] //controllo se nella mappa c'è la chiave message.TS

	if ok == true {
		//fmt.Println("\nla chiave e' presente -->", map1[message.SeqNum[0]])
		//prendo elementi value relativi a quella chiave e faccio controllo
		if len(map1[message.TS]) == 1 { //c'e un solo valore nella lista relativa al TS (1 solo msg con quel TS)
			if message.Sender < map1[message.TS][0].Sender {
				map1[message.TS] = append([]lamport.Message{message}, map1[message.TS]...) //inserisco il msg all'inizio dello slice
				/*
					NB: data := []string{"A", "B", "C", "D"}

						METTO ALLA FINE: data = append(data, "prova")	--> [A B C D prova]

						METTO ALL'INIZIO: data = append([]string{"prova"}, data...) --> [prova A B C D]
				*/
			}
		} else {
			//fmt.Println("sono in else")
			for i := 1; i < len(map1[message.TS]); i++ {
				if map1[message.TS][i-1].Sender < message.Sender && message.Sender < map1[message.TS][i].Sender {
					//fmt.Println("IL MSG STA TRA I 2")
					// devo inserire il msg tra i-1 e i
					map1[message.TS] = append(map1[message.TS], message) // msg ora e' in posiz len(msg1)
					//fmt.Println("in else map1 ===", map1)
					//fmt.Println("i ===", i)
					//fmt.Println("fine ==", map1[message.SeqNum[0]][len(map1[message.SeqNum[0]])-1])
					//map1[message.TS][i], map1[message.TS][len(map1[message.TS])-1] = map1[message.TS][len(map1[message.TS])-1], map1[message.TS][i]
					copy(map1[message.TS][i+1:], map1[message.TS][i:])
					//fmt.Println("dopo copy map1 ===", map1)
					map1[message.TS][i] = message
					/*
						se ho slice : arr =[1 3 5] e voglio aggiungere il 2 tra 1 e 3:
						1. metto il 2 alla fine --> arr = [1 3 5 2]
						2. arr[2:] == [5 2] e arr[1:] == [3 5 2], con copy copio [3 5 2] in [5 2], ottenendo arr = [1 3 3 5]
						3. poi dico che arr[1] = 2 --> arr = [1 2 3 5]
					*/

					break
				}
				if map1[message.TS][len(map1[message.TS])-1].Sender < message.Sender {
					map1[message.TS] = append(map1[message.TS], message) // metto msg alla fine
					break
				}
			}
		}
		/*
			for i := 0; i < len(map1[message.TS]); i++ {
				fmt.Println("boh --- ", map1[message.TS][i])
				if message.Sender < map1[message.TS][i].Sender {
					map1[message.TS] = append(map1[message.TS], message)
					//metto booleano per dire se c'è o no
				}
			}

		*/
		//map1[message.TS] = append(map1[message.TS], message)

	} else { // nella mappa non c'e quella chiave
		//fmt.Println("la chiave non e' presente ")
		map1[message.TS] = append(listMsg, message)
	}
	//map1[message2.TS] = append(listMsg, message2)
	fmt.Println("map1 ===", map1)
}
func GetFirstElementMap(mapMsg MessageMap) lamport.Message {
	var message lamport.Message
	for key, element := range mapMsg {
		fmt.Println("Key:", key, "=>", "Element:", element)
		message = element[0]
		break
	}
	fmt.Println("GetFirstElementMap ------", message)
	return message

}

func RemoveFirstElementMap(mapMsg MessageMap) {

	for key, _ := range mapMsg {
		mapMsg[key] = mapMsg[key][1:]
		if len(mapMsg[key]) == 0 { //se non ci sono piu msg con quel TS, ossia la lista di msg per quel TS (key) e' vuota
			delete(mapMsg, key)

			break
		}
	}

	fmt.Println("mappa == ", mapMsg)

}
