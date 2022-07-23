package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeMap(t *testing.T) {
	out := make(chan []int, 2)
	Apply(
		Of(1, 2),
		Pipe2(
			MergeMap(func(v int) Observable[[]int] {
				return Apply(
					Range(0, v),
					ToSlice[int](),
				)
			}),
			ToChan(out),
		),
	).Subscribe(NewSubscriberEmpty[[]int]())

	expext := [][]int{{0, 1}, {0, 1, 2}}
	for o := range out {
		require.Contains(t, expext, o)
	}
}
