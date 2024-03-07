package tests

import (
	"challenge/pkg/proto"
	suits "challenge/tests/suit"
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
)

func TestStartTimer_Ok(t *testing.T) {
	_, s := suits.NewDefault(t)

	timerName := gofakeit.Word()
	var freq int64 = 3
	var secs int64 = 14

	// Start the connection which will create timer
	c, err := s.Client.StartTimer(context.Background(), &proto.Timer{Name: timerName, Seconds: secs, Frequency: freq})
	require.NoError(t, err)

	// Start tracking timer updates (small 1 second delta applied to exclude server delays)
	expectedSecs := secs - freq
	for {
		timer, err := c.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		assert.Equal(t, timerName, timer.GetName())
		assert.Equal(t, freq, timer.GetFrequency())
		// Small one second delta, because server, due to api delays, may send the same second or next second with 1 step difference
		if timer.GetSeconds() != expectedSecs && timer.GetSeconds() != expectedSecs-1 {
			t.Fail()
		}
		expectedSecs -= freq
	}
}

func TestStartTimer_OkWithReconnect(t *testing.T) {
	_, s := suits.NewDefault(t)

	timerName := gofakeit.Word()
	var freq int64 = 1
	var secs int64 = 10
	var delay int64 = 5

	// Start first connection, which will create our timer
	c, err := s.Client.StartTimer(context.Background(), &proto.Timer{Name: timerName, Seconds: secs, Frequency: freq})
	require.NoError(t, err)

	// Immediately close connection, timer must keep running
	err = c.CloseSend()
	require.NoError(t, err)

	// Wait specific amount of time
	time.Sleep(time.Duration(delay) * time.Second)

	// Reconnect to the same timer name, we should get the same timer we created before
	c, err = s.Client.StartTimer(context.Background(), &proto.Timer{Name: timerName, Seconds: secs, Frequency: freq})
	require.NoError(t, err)

	// Start tracking timer updates like we were not closing the connection
	expectedSecs := secs - delay
	for {
		timer, err := c.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		assert.Equal(t, timerName, timer.GetName())
		assert.Equal(t, freq, timer.GetFrequency())
		// Small one second delta, because server, due to api delays, may send the same second or next second with 1 step difference
		if timer.GetSeconds() != expectedSecs && timer.GetSeconds() != expectedSecs-1 {
			t.Fatalf("expected: %d, got: %d", expectedSecs, timer.GetSeconds())
		}
		expectedSecs -= freq
	}
}

func TestStartTimer_OkWithDifferentClients(t *testing.T) {
	_, s := suits.NewDefault(t)

	timerName := gofakeit.Word()
	var secs int64 = 10
	// Every client has it's own frequency value
	// That means they receive updates independently
	var freq1 int64 = 1
	var freq2 int64 = 3

	// Start first connection, with the first client
	c1, err := s.Client.StartTimer(context.Background(), &proto.Timer{Name: timerName, Seconds: secs, Frequency: freq1})
	require.NoError(t, err)

	// Start second connection, with the second client
	c2, err := s.Client.StartTimer(context.Background(), &proto.Timer{Name: timerName, Seconds: secs, Frequency: freq2})
	require.NoError(t, err)

	// We start 2 independent linear loops because we dont care of the time when our message arrives
	// We have to make sure that we have all needed messages for each client
	exp1 := secs - freq1
	for {
		timer1, err1 := c1.Recv()
		if err1 == io.EOF {
			break
		}
		require.NoError(t, err1)

		assert.Equal(t, timerName, timer1.GetName())
		assert.Equal(t, freq1, timer1.GetFrequency())
		// Small one second delta, because server, due to api delays, may send the same second or next second with 1 step difference
		if timer1.GetSeconds() != exp1 && timer1.GetSeconds() != exp1-1 {
			t.Fatalf("expected: %d, got: %d", exp1, timer1.GetSeconds())
		}
		exp1 -= freq1
	}

	exp2 := secs - freq2
	for {
		timer2, err2 := c2.Recv()
		if err2 == io.EOF {
			break
		}
		require.NoError(t, err2)

		assert.Equal(t, timerName, timer2.GetName())
		assert.Equal(t, freq2, timer2.GetFrequency())
		// Small one second delta, because server, due to api delays, may send the same second or next second with 1 step difference
		if timer2.GetSeconds() != exp2 && timer2.GetSeconds() != exp2-1 {
			t.Fatalf("expected: %d, got: %d", exp2, timer2.GetSeconds())
		}
		exp2 -= freq2
	}
}
