package rx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMergeAll(t *testing.T) {
	out := make(chan []int, 1)

	Apply(
		Of(
			Of(1),
			Of(2),
			Range(3, 4),
			Apply(
				Interval(10*time.Millisecond),
				Pipe2(Take[int](2), Map(func(v int) (int, error) { return v + 10, nil })),
			),
		),
		Pipe3(
			MergeAll[int](),
			ToSlice[int](),
			ToChan(out),
		),
	).Subscribe(NewSubscriberEmpty[[]int]())

	res := []int{1, 2, 3, 4, 10, 11}
	v := <-out

	for _, vv := range v {
		require.Contains(t, res, vv)
	}
}
