package timer

type SubUnsub struct {
	timers map[string]map[chan Ping]bool
}

func NewSubUnsub() *SubUnsub {
	return &SubUnsub{
		timers: make(map[string]map[chan Ping]bool),
	}
}

func (t *SubUnsub) Sub(timerName string, c chan Ping) {
	if t.timers[timerName] == nil {
		t.timers[timerName] = make(map[chan Ping]bool)
	}
	t.timers[timerName][c] = true
}

func (t *SubUnsub) UnsubAll(timerName string) {
	for c, _ := range t.timers[timerName] {
		close(c)
	}
	delete(t.timers, timerName)
}

func (t *SubUnsub) Unsub(timerName string, c chan Ping) {
	if t.timers[timerName] == nil {
		return
	}

	close(c)
	delete(t.timers[timerName], c)
}
