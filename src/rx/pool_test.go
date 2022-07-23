package rx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPoolInfinite(t *testing.T) {
	out := make(chan []int)
	go Apply(
		Range(0, 5),
		Pipe3(
			PoolInfinite(
				Map(func(v int) (int, error) {
					time.Sleep(time.Duration(250-v*25) * time.Millisecond)
					return v, nil
				}),
			),
			ToSlice[int](),
			ToChan(out),
		),
	).Subscribe(
		NewSubscriberEmpty[[]int](),
	)

	v := <-out
	require.Equal(t, []int{5, 4, 3, 2, 1, 0}, v)
}

func TestPool(t *testing.T) {
	out := make(chan []int)
	go Apply(
		Range(0, 5),
		Pipe3(
			Pool(
				6,
				Map(func(v int) (int, error) {
					time.Sleep(time.Duration(250-v*25) * time.Millisecond)
					return v, nil
				}),
			),
			ToSlice[int](),
			ToChan(out),
		),
	).Subscribe(
		NewSubscriberEmpty[[]int](),
	)
	v := <-out
	require.Equal(t, 6, len(v))
	require.Equal(t, []int{5, 4, 3, 2, 1, 0}, v)
}
