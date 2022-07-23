package rx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntervalUncomplete(t *testing.T) {
	un := Interval(1 * time.Millisecond).
		Subscribe(NewSubscriberFn(
			func(v int) {},
			func(e error) { t.Error("should not error") },
			func() { t.Error("should not complete") },
		))
	un()
	time.Sleep(1 * time.Millisecond)
}

func TestIntervalComplete(t *testing.T) {
	d := make(chan any)
	c := 0
	Apply(Interval(1*time.Millisecond), Take[int](2)).
		Subscribe(NewSubscriberFn(
			func(v int) {},
			func(e error) { c++ },
			func() {
				c++
				d <- true
			},
		))
	<-d
	require.Equal(t, 1, c)
}
