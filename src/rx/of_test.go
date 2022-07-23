package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRange(t *testing.T) {
	var r []int
	d := make(chan any, 1)

	Apply(
		Range(0, 9999999999),
		Pipe4(
			Take[int](3),
			ToSlice[int](),
			Do(func(v []int) {
				r = v
			}),
			Finalize[[]int](func() {
				close(d)
			}),
		),
	).Subscribe(NewSubscriberEmpty[[]int]())

	<-d

	require.Equal(t, []int{0, 1, 2}, r)
}
