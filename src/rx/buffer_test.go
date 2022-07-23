package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {
	cout := make(chan []int)
	Apply(
		Range(1, 6),
		Pipe2(
			Buffer[int](2),
			ToChan(cout),
		),
	).Subscribe(NewSubscriberEmpty[[]int]())

	out := [][]int{}
	for v := range cout {
		out = append(out, v)
	}
	require.Equal(t, [][]int{
		{1, 2},
		{3, 4},
		{5, 6},
	}, out)
}
