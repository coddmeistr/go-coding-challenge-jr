package timercheck

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

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
				timerName: "test",
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
				timerName: "SomeOtherName",
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
				timerName: "timeout",
				timerTime: 4,
				testWait:  5,
			},
			wantRemain:  0,
			wantElapsed: 0,
			wantError:   true,
			wantErrMsg:  "timer timed out",
		},
	}

	timercheck := NewTimerCheck(http.DefaultClient)
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := timercheck.CreateTimer(tt.args.timerName, tt.args.timerTime)
			require.NoError(t, err)

			time.Sleep(time.Duration(tt.args.testWait) * time.Second)

			remain, elapsed, err := timercheck.CheckTimer(tt.args.timerName)
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.Equal(t, tt.wantRemain, remain)
				assert.Equal(t, tt.wantElapsed, elapsed)
			}

		})
	}
}
