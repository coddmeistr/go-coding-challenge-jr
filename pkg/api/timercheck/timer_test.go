package timercheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func newClientMock(t *testing.T, statusCode int, path string, response any) *TimerCheck {
	return &TimerCheck{
		client: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, path, r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				respBody, err := json.Marshal(response)
				if err != nil {
					assert.Fail(t, "Cannot read bytes")
				}
				return &http.Response{
					StatusCode: statusCode,
					Body:       io.NopCloser(bytes.NewReader(respBody)),
				}, nil
			}),
		},
	}
}

func TestCheckTimer_TestCases(t *testing.T) {

	tc := []struct {
		name string

		timerName string

		expectedStatusCode int
		expectedResponse   any
		wantErr            bool
		wantErrMsg         string
	}{
		{
			name:               "ok",
			timerName:          "test",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   TimerResponse{Remaining: 10, Elapsed: 2},
			wantErr:            false,
			wantErrMsg:         "",
		},
		{
			name:               "some error",
			timerName:          "test",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   TimerResponse{Remaining: 0, Elapsed: 0},
			wantErr:            true,
			wantErrMsg:         "bad http status code",
		},
		{
			name:               "when timed out",
			timerName:          "test",
			expectedStatusCode: 504,
			expectedResponse:   TimerResponse{},
			wantErr:            true,
			wantErrMsg:         "timed out",
		},
		{
			name:               "when not exists",
			timerName:          "test",
			expectedStatusCode: 404,
			expectedResponse:   TimerResponse{},
			wantErr:            true,
			wantErrMsg:         "not exists",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			timer := newClientMock(t, tt.expectedStatusCode,
				fmt.Sprintf("/%s", tt.timerName),
				tt.expectedResponse)

			rem, el, err := timer.CheckTimer(tt.timerName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
				resp, ok := tt.expectedResponse.(TimerResponse)
				require.True(t, ok)
				assert.Equal(t, int(resp.Remaining), rem)
				assert.Equal(t, int(resp.Elapsed), el)
			}
		})
	}
}

func TestCreateTimer_TestCases(t *testing.T) {

	type args struct {
		timerName string
		timerSecs int
	}

	tc := []struct {
		name string

		args args

		expectedStatusCode int
		expectedResponse   any
		wantErr            bool
		wantErrMsg         string
	}{
		{
			name: "ok",
			args: args{
				timerName: "test",
				timerSecs: 10,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   TimerResponse{Remaining: 10, Elapsed: 0},
			wantErr:            false,
			wantErrMsg:         "",
		},
		{
			name: "some error",
			args: args{
				timerName: "test",
				timerSecs: 10,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   TimerResponse{Remaining: 0, Elapsed: 0},
			wantErr:            true,
			wantErrMsg:         "bad http status code",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			timer := newClientMock(t, tt.expectedStatusCode,
				fmt.Sprintf("/%s/%d", tt.args.timerName, tt.args.timerSecs),
				tt.expectedResponse)

			err := timer.CreateTimer(tt.args.timerName, tt.args.timerSecs)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// API test
func TestTimerCheck_TestCases(t *testing.T) {
	type args struct {
		timerName string
		timerTime int
		testWait  int
	}

	tc := []struct {
		name        string
		args        args
		wantRemain  int
		wantElapsed int
		wantError   bool
		wantErrMsg  string
	}{
		{
			name: "ok",
			args: args{
				timerName: gofakeit.Username(),
				timerTime: 3,
				testWait:  2,
			},
			wantRemain:  1,
			wantElapsed: 2,
			wantError:   false,
			wantErrMsg:  "",
		},
		{
			name: "ok, diff name",
			args: args{
				timerName: gofakeit.Username(),
				timerTime: 2,
				testWait:  1,
			},
			wantRemain:  1,
			wantElapsed: 1,
			wantError:   false,
			wantErrMsg:  "",
		},
		{
			name: "timeout",
			args: args{
				timerName: gofakeit.Username(),
				timerTime: 4,
				testWait:  5,
			},
			wantRemain:  0,
			wantElapsed: 0,
			wantError:   true,
			wantErrMsg:  "timer timed out",
		},
	}

	timerCheck := NewTimerCheck(http.DefaultClient)
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := timerCheck.CreateTimer(tt.args.timerName, tt.args.timerTime)
			require.NoError(t, err)

			// waiting some amount of time but applying small 1/2 second delta
			// we're doing it to exclude server delays
			delta := time.Duration(200) * time.Millisecond
			time.Sleep(time.Duration(tt.args.testWait)*time.Second - delta)

			remain, elapsed, err := timerCheck.CheckTimer(tt.args.timerName)
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantRemain, remain)
				assert.Equal(t, tt.wantElapsed, elapsed)
			}

		})
	}
}

// API test
func TestTimerCheck_NotExistingTimer(t *testing.T) {
	timerCheck := NewTimerCheck(http.DefaultClient)
	notExistingTimerName := gofakeit.Username() // Unique timer name
	_, _, err := timerCheck.CheckTimer(notExistingTimerName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not exists")
}
