package timer

import (
	"challenge/pkg/api/timercheck"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Ping struct {
	TimerName   string
	SecondsLeft int
}

type Timer struct {
	timerChecker timercheck.TimerCheck
	su           *SubUnsub
}

func NewTimer(timerChecker timercheck.TimerCheck) *Timer {
	return &Timer{
		timerChecker: timerChecker,
		su:           NewSubUnsub(),
	}
}

// StartOrSubscribe creating new streaming channel which gets timer updates with given frequency
//
//		Streaming was created only if there was no errors in return
//
//	 You can use context cancellation function to interrupt streaming from outside
//
// NOTE: Subscribe recommended to use instead, because it reduces API calls due to broadcasting system
func (t *Timer) StartOrSubscribe(timerName string, timerSeconds int, freq int) (<-chan Ping, context.CancelFunc, error) {

	_, _, err := t.timerChecker.CheckTimer(timerName)
	if err != nil {
		if errors.Is(err, timercheck.ErrTimedOut) || errors.Is(err, timercheck.ErrNotExists) {
			log.Println("timer doesn't exist, creating new timer with name: " + timerName)
			if err := t.timerChecker.CreateTimer(timerName, timerSeconds); err != nil {
				return nil, nil, fmt.Errorf("%w: %v", err, "timer creation failed")
			}
		} else {
			return nil, nil, fmt.Errorf("%w: %v", err, "something wrong on foreign api side")
		}
	}

	ticker := time.NewTicker(time.Duration(freq) * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	ping := make(chan Ping)
	go func() {
		defer func() {
			close(ping)
			ticker.Stop()
			log.Println("returning from StartOrSubscribe timer goroutine")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r, _, err := t.timerChecker.CheckTimer(timerName)
				if err != nil {
					if errors.Is(err, timercheck.ErrTimedOut) {
						return
					}
					log.Println("error when checking timer: ", err)
					return
				}
				ping <- Ping{
					TimerName:   timerName,
					SecondsLeft: r,
				}
			}
		}

	}()

	return ping, cancel, nil
}

// Subscribe subscribes to timer updates on returned channel
//
// If some timer currently running, it will return streaming channel bound to frequency
// that was set when timer was firstly created
//
//	If timer not exists or timed out, it will create new broadcast goroutine and subscribe new channel to this goroutine
//
// When timer expires, all subscribed channels will be automatically unsubscribed(closed)
func (t *Timer) Subscribe(timerName string, timerSeconds int, freq int) (chan Ping, error) {

	_, _, err := t.timerChecker.CheckTimer(timerName)
	if err == nil {
		// If timer already running, subscribe to it
		log.Println("timer already running with name: " + timerName)
		log.Println("update frequency binds may be different")
		c := make(chan Ping)
		t.su.Sub(timerName, c)
		return c, nil
	} else {
		if !errors.Is(err, timercheck.ErrTimedOut) && !errors.Is(err, timercheck.ErrNotExists) {
			log.Println("error when checking timer: " + timerName)
			return nil, fmt.Errorf("%w: %v", err, "timer creation failed")
		}
	}

	// Create timer and subscribe new channel
	if err := t.timerChecker.CreateTimer(timerName, timerSeconds); err != nil {
		return nil, fmt.Errorf("%w: %v", err, "timer creation failed")
	}

	c := make(chan Ping)
	t.su.Sub(timerName, c)

	// Create broadcasting goroutine
	ticker := time.NewTicker(time.Duration(freq) * time.Second)
	go func() {
		defer func() {
			t.su.UnsubAll(timerName)
			ticker.Stop()
			log.Println("returning from Subscribe timer goroutine")
		}()

		for {
			select {
			case <-ticker.C:
				r, _, err := t.timerChecker.CheckTimer(timerName)
				if err != nil {
					if errors.Is(err, timercheck.ErrTimedOut) {
						return
					}
					log.Println("error when checking timer: ", err)
					return
				}
				p := Ping{
					TimerName:   timerName,
					SecondsLeft: r,
				}
				log.Printf("streaming to %d subscribed channels\n", len(t.su.timers[timerName]))
				for c, _ := range t.su.timers[timerName] {
					c <- p
				}
			}
		}

	}()

	return c, nil
}

// Unsubscribe simply deletes given channel bound to given timer name
// from broadcast system
func (t *Timer) Unsubscribe(timerName string, c chan Ping) {
	t.su.Unsub(timerName, c)
}
