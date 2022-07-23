package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTake(t *testing.T) {
	complete := 0
	var r []int

	d, f := FinalizeCreateDone[int]()

	Apply(
		Of(1, 2, 3, 4, 5),
		Pipe3(
			Take[int](2),
			f,
			ToSlice[int](),
		),
	).Subscribe(NewSubscriberFn(
		func(v []int) { r = v },
		func(e error) {
			t.Error("had no error")
		},
		func() {
			complete++
		},
	))
	<-d
	require.Equal(t, []int{1, 2}, r)
	require.Equal(t, 1, complete)
}
