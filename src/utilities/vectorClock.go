package utilities

import (
	"fmt"
	"sort"
	"strconv"
)

type VectorClock = map[string]int

func StartVC2(vc VectorClock) {
	numKeys := MAXPEERS - 1
	for i := 0; i < numKeys; i++ {
		username := "p" + strconv.Itoa(i+1)
		vc[username] = 0
	}
}

func ToString(vc VectorClock) string {

	keys := make([]string, 0, len(vc))
	values := make([]int, 0, len(vc))

	for k, _ := range vc {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, vc[k])
		values = append(values, vc[k])
	}
	vcString := fmt.Sprint(values)
	return vcString
}

func IsEligible(vc0 VectorClock, vcPeer VectorClock, usernamePeer string) bool {
	var eligible bool
	for username, clock := range vc0 {
		if username != usernamePeer {
			if vcPeer[username] > clock {
				eligible = false
				break
			} else {
				eligible = true
			}
		}
	}
	return eligible

}

func IncrementVC(vc VectorClock, username string) {
	vc[username] = vc[username] + 1
}

func UpdateVC(vc VectorClock, vcMsg VectorClock) {
	for k, _ := range vc {
		max := max(vc[k], vcMsg[k])
		vc[k] = max
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
