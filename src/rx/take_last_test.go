package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTakeLast(t *testing.T) {
	done := make(chan any, 1)
	var r []int
	Apply(
		Range(1, 10),
		TakeLast[int](2),
	).Subscribe(NewSubscriberFn(
		func(v []int) {
			r = v
		},
		func(e error) {},
		func() {
			close(done)
		},
	))
	<-done
	require.Equal(t, []int{9, 10}, r)
}
