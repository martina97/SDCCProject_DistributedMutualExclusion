package utilities

import (
	"fmt"
	"strconv"
)

type VectorClock = map[string]int

func ToString(vc VectorClock) string {
	values := make([]int, 0, len(vc))

	for _, v := range vc {
		values = append(values, v)
	}
	fmt.Println("values =", values)
	string := fmt.Sprint(values)
	return string
}

/*
func (vc *VectorClock) ToString2() string {
	values := make([]int, 0, len(*vc))

	for _, v := range *vc {
		values = append(values, v)
	}
	fmt.Println("values =", values)
	string := fmt.Sprint(values)
	return string
}

*/

func StartVC2(vc VectorClock) {
	//vc = make(map[string]int)
	fmt.Println("sto in StartVC")
	numKeys := MAXPEERS - 1
	for i := 0; i < numKeys; i++ {
		username := "p" + strconv.Itoa(i+1)
		fmt.Println("username == ", username)
		vc[username] = 0
	}
}

func IsEligible(vc0 VectorClock, vcPeer VectorClock, usernamePeer string) bool {
	var eligible bool
	for username, clock := range vc0 {
		if username != usernamePeer {
			fmt.Println("VC0[", username, "] =", clock)
			fmt.Println("VC2[", username, "] =", vcPeer[username])
			if vcPeer[username] > clock {
				eligible = false
				break
			} else {
				eligible = true
			}
		}
	}
	fmt.Println("eligible == ", eligible)
	return eligible

}

func IncrementVC(vc VectorClock, username string) {

	vc[username] = vc[username] + 1

}
