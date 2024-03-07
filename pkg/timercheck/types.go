package timercheck

type TimerResponse struct {
	Elapsed   float64 `json:"seconds_elapsed"`
	Remaining float64 `json:"seconds_remaining"`
}
