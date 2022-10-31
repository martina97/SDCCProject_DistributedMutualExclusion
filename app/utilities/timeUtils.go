package utilities

const maxDelay = 5

/*
type Clock interface {
	Start()           // Start not concurrent operation
	Increment(id int) //id: my node identifier, do not care if ScalarClock
	Update(timestamp []uint64)
	GetValue() []uint64
	//AtomicIncAndGet ??
}

type ScalarClock struct {
	counter uint64 //todo : ossia valore TS ?!
	mutex   sync.Mutex
}

func (clock *ScalarClock) Start() {
	clock.counter = 0
}

func (clock *ScalarClock) Increment(_id int) {
	//serve sia quando mando sia quando ricevo msg
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter++
}
func (clock *ScalarClock) Update(timestamp []uint64) {
	//serve quando ricevo msg
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	clock.counter = MaxOf(clock.counter, timestamp[0])
}
func (clock *ScalarClock) GetValue() []uint64 {
	ret := make([]uint64, 1)
	clock.mutex.Lock()
	defer clock.mutex.Unlock()
	ret[0] = clock.counter
	return ret
}

func Delay() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(maxDelay)
	time.Sleep(time.Duration(n) * time.Second)
}

func Delay_ms(maxTime int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(maxTime)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func Delay_sec(exactTime int) {
	time.Sleep(time.Duration(exactTime) * time.Second)
}

func Timer(timeout int, outChan chan bool) {
	Delay_sec(timeout)
	outChan <- true
}

func MaxOf(vars ...uint64) uint64 {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}


*/
