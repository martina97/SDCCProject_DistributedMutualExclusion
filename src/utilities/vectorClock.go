package utilities

import "fmt"

type VectorClock [MAXPEERS]int //todo: invece di mettere 3, devo decidere il # peer in modo automatico

func StartVC(vc VectorClock) {
	fmt.Println("sto in StartVC")
	for i := 0; i < len(vc); i++ {
		fmt.Println(vc[i])
		vc[i] = 0
		fmt.Println(vc[i])
	}
}
