package utilities

type TimeStamp int

func StartTS(ts TimeStamp) {
	ts = 0
}

func IncrementTS(ts *TimeStamp) {
	*ts++
}

func UpdateTS(ts, tsMsg *TimeStamp, algo string) {
	switch algo {
	case "Lamport":
		*ts = Max(*ts, *tsMsg) + 1 // se Lamport, ts = max(ts, tsMessaggio) + 1
	case "ricartAgrawala":
		*ts = Max(*ts, *tsMsg) //se ricartAgrawala, Num = max(Num, tsMessaggio)
	}

}

func Max(x, y TimeStamp) TimeStamp {
	if x < y {
		return y
	}
	return x
}
