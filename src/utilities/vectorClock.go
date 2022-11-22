package utilities

import (
	"fmt"
	"sort"
	"strconv"
)

type VectorClock = map[string]int

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
func ToString(vc VectorClock) string {
	fmt.Println("sto in ToString dentro VectorClock, vc == ", vc)
	keys := make([]string, 0, len(vc))
	values := make([]int, 0, len(vc))

	for k, v := range vc {
		fmt.Println("k = ", k)
		fmt.Println("v = ", v)
		//values = append(values, v)
		keys = append(keys, k)
	}
	fmt.Println("keys =", keys)
	sort.Strings(keys)
	fmt.Println("keys =", keys)
	for _, k := range keys {
		fmt.Println(k, vc[k])
		values = append(values, vc[k])
	}
	fmt.Println("values =", values)

	string := fmt.Sprint(values)

	return string
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

func UpdateVC(vc VectorClock, vcMsg VectorClock) {
	for k, _ := range vc {
		max := max(vc[k], vcMsg[k])
		fmt.Println("max == ", max)
		vcMsg[k] = max
	}
	fmt.Println("vcMsg ==", vcMsg)

}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}