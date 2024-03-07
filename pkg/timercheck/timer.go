package timercheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrInternal  = errors.New("internal library error")
	ErrTimedOut  = errors.New("timer timed out")
	ErrNotExists = errors.New("timer not exists")
)

const (
	host = "https://timercheck.io/"
)

type TimerCheck struct {
	client *http.Client
}

func NewTimerCheck(c *http.Client) *TimerCheck {
	return &TimerCheck{
		client: c,
	}
}

// CreateTimer creates new timer using timercheck.io API
// It created new timer using provided name and with timer seconds of provided value
func (t *TimerCheck) CreateTimer(name string, seconds int) error {

	req, err := http.NewRequest("GET", host+name+"/"+fmt.Sprintf("%d", seconds), nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%w: %v", ErrInternal, "got bad http status code")
	}

	return nil
}

func (t *TimerCheck) CheckTimer(name string) (remain int, elapsed int, err error) {

	req, err := http.NewRequest("GET", host+name, nil)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrInternal, err)
		return
	}

	resp, err := t.client.Do(req)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrInternal, err)
		return
	}

	if resp.StatusCode == 504 {
		err = fmt.Errorf("%w: %v", ErrTimedOut, "timer timed out")
		return
	}

	if resp.StatusCode == 404 {
		err = fmt.Errorf("%w: %v", ErrNotExists, "timer never been created")
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%w: %v", ErrInternal, "got bad http status code")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrInternal, err)
		return
	}

	var timerResp TimerResponse
	err = json.Unmarshal(body, &timerResp)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrInternal, err)
		return
	}

	remain = int(timerResp.Remaining)
	elapsed = int(timerResp.Elapsed)
	return
}
