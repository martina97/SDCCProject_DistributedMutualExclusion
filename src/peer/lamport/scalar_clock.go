package lamport

type ScalarClock int

func StartSC(ts ScalarClock) {
	ts = 0
}

func IncrementSC(ts *ScalarClock) {
	*ts++
}

func UpdateSC(ts, tsMsg *ScalarClock) {
	*ts = Max(*ts, *tsMsg)
}

func Max(x, y ScalarClock) ScalarClock {
	if x < y {
		return y
	}
	return x
}
