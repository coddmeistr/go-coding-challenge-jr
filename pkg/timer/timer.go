package timer

import (
	"challenge/pkg/timercheck"
	"errors"
	"fmt"
	"time"
)

type Ping struct {
	TimerName   string
	SecondsLeft int
}

type Timer struct {
	timerChecker timercheck.TimerCheck
}

func NewTimer(timerChecker timercheck.TimerCheck) *Timer {
	return &Timer{
		timerChecker: timerChecker,
	}
}

func (t *Timer) StartOrSubscribe(timerName string, timerSeconds int, freq int) (<-chan Ping, chan<- struct{}, error) {

	// Error timed out from timer checker considering that the
	// timer is not created or expired
	_, _, err := t.timerChecker.CheckTimer(timerName)
	if err != nil {
		if errors.Is(err, timercheck.ErrTimedOut) {
			if err := t.timerChecker.CreateTimer(timerName, timerSeconds); err != nil {
				return nil, nil, fmt.Errorf("%w: %v", err, "timer creation failed")
			}
		}

		return nil, nil, fmt.Errorf("%w: %v", err, "something wrong on foreign api side")
	}

	ticker := time.NewTicker(time.Duration(freq) * time.Second)
	done := make(chan struct{})
	ping := make(chan Ping)
	go func() {
		defer close(ping)
		defer fmt.Println("Closing timer timer subscription goroutine")

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				r, _, err := t.timerChecker.CheckTimer(timerName)
				if err != nil {
					return
				}
				ping <- Ping{
					TimerName:   timerName,
					SecondsLeft: r,
				}
			}
		}

	}()

	return ping, done, nil
}
