package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	var r []int
	d := make(chan any)

	Apply(
		Of(1, 2, 3),
		Pipe2(
			Scan(func(acc int, v int) (int, error) { return acc + v, nil }, 1),
			ToSlice[int](),
		),
	).Subscribe(NewSubscriberFn(
		func(v []int) { r = v },
		func(e error) { t.Error("should not error") },
		func() { d <- true },
	))

	<-d
	require.Equal(t, []int{2, 4, 7}, r)
}
