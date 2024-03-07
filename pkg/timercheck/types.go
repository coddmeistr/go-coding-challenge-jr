package timercheck

type TimerResponse struct {
	Elapsed   int `json:"seconds_elapsed"`
	Remaining int `json:"seconds_remaining"`
}
