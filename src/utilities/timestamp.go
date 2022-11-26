package utilities

type TimeStamp int

func StartTS(ts TimeStamp) {
	ts = 0
}

func IncrementTS(ts *TimeStamp) {
	*ts++
}

func UpdateTS(ts, tsMsg *TimeStamp) {
	*ts = Max(*ts, *tsMsg)
}

func Max(x, y TimeStamp) TimeStamp {
	if x < y {
		return y
	}
	return x
}
