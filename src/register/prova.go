package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

var vectorClock utilities.VectorClock

func main2() {

	VC0 := make(map[string]int)
	VC0["p1"] = 0
	VC0["p2"] = 0
	VC0["p3"] = 0
	fmt.Println(VC0)
	//value := VC0["p1"]
	//fmt.Println(value)

	VC2 := make(map[string]int)
	VC2["p1"] = 1
	VC2["p2"] = 1
	VC2["p3"] = 0
	fmt.Println(VC2)
	//username := "p2"
	//devo controllare che p2 puo avere il token, quindi controllo se
	//il valore del clock di p1 in VC2 <= valore del clock di p1 in VC0
	var eligible bool
	for username, clock := range VC0 {
		if username != "p2" {
			fmt.Println("VC0[", username, "] =", clock)
			fmt.Println("VC2[", username, "] =", VC2[username])
			if VC2[username] > clock {
				eligible = false
				break
			} else {
				eligible = true
			}
		}
	}
	fmt.Println("eligible == ", eligible)

	vectorClock := make(map[string]int)
	utilities.StartVC2(vectorClock)
	fmt.Println("vectorClock == ", vectorClock)

	utilities.IsEligible(VC0, VC2, "p2")

}
