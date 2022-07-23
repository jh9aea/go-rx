package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToSlice(t *testing.T) {
	out := make(chan int)
	Apply(
		Of(1, 2, 3),
		ToChan(out),
	).Subscribe(NewSubscriberEmpty[int]())

	r := []int{1, 2, 3}
	i := 0
	for v := range out {
		require.Equal(t, r[i], v)
		i++
	}
}
