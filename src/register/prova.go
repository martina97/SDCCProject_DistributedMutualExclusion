package main

import (
	"SDCCProject_DistributedMutualExclusion/src/utilities"
	"fmt"
)

var vectorClock utilities.VectorClock

func main2() {

	enterString1 := "[10:56:27.833] : p1 enters the critical section at 10:56:27.833."

	exitString1 := "[10:56:57.863] : p1 exits the critical section at 10:56:57.863."
	enterString2 := "[10:55:57.726] : p2 enters the critical section at 10:55:57.726."
	exitString2 := "[10:56:27.758] : p2 exits the critical section at 10:56:27.758."

	enterCS := "enters the critical section at "
	exitCS := "exits the critical section at "

	/*

		dateStringEnter1 := after(enterString1, enterCS)
		fmt.Println("dateStringEnter1 ==", dateStringEnter1)

		dateStringExit1 := after(exitString1, exitCS)
		fmt.Println("dateStringExit1 ==", dateStringExit1)

		dateStringEnter2 := after(enterString2, enterCS)
		fmt.Println("dateStringEnter2 ==", dateStringEnter2)

		dateStringExit2 := after(exitString2, exitCS)
		fmt.Println("dateStringExit2 ==", dateStringExit2)

		dateStringExit1 = strings.ReplaceAll(dateStringExit1, ":", "")
		dateStringEnter2 = strings.ReplaceAll(dateStringEnter2, ":", "")
		dateStringExit2 = strings.ReplaceAll(dateStringExit2, ":", "")
		dateStringEnter1 = strings.ReplaceAll(dateStringEnter1, ":", "")

		enter1, _ := time.Parse("150405.000.", dateStringEnter1)
		exit1, _ := time.Parse("150405.000.", dateStringExit1)
		enter2, _ := time.Parse("150405.000.", dateStringEnter2)
		exit2, _ := time.Parse("150405.000.", dateStringExit2)
		fmt.Println("enter1 ==", enter1)
		fmt.Println("exit1 ==", exit1)
		fmt.Println("enter2 ==", enter2)
		fmt.Println("exit2 ==", exit2)
		fmt.Println("dateStringEnter1 ==", dateStringEnter1)

		fmt.Println("dateStringExit1 ==", dateStringExit1)

		fmt.Println("dateStringEnter2 ==", dateStringEnter2)

		fmt.Println("dateStringExit2 ==", dateStringExit2)

	*/

	enter1 := utilities.ConvertStringToDate(enterString1, enterCS)
	exit1 := utilities.ConvertStringToDate(exitString1, exitCS)
	enter2 := utilities.ConvertStringToDate(enterString2, enterCS)
	exit2 := utilities.ConvertStringToDate(exitString2, exitCS)
	fmt.Println("enter1 ==", enter1)
	fmt.Println("exit1 ==", exit1)
	fmt.Println("enter2 ==", enter2)
	fmt.Println("exit2 ==", exit2)

	//fmt.Println("enter1.Before(enter2) ?", enter1.Before(enter2))
	if enter1.Before(enter2) {
		fmt.Println("ENTRA PRIMA p1")
		//vuol dire che p1 entra prima di p2, quindi devo controllare che exit1<enter2
		fmt.Println("exit1.Before(enter2) ?", exit1.Before(enter2))
	} else {
		fmt.Println("ENTRA PRIMA p2")
		fmt.Println("exit2 ==", exit2)
		fmt.Println("enter1 ==", enter1)
		fmt.Println("exit2.Before(enter1) ?", exit2.Before(enter1))

	}

}
