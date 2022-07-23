package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupBy(t *testing.T) {
	type V struct {
		Id   int
		Name string
	}
	out := make(chan []V, 3)
	out2 := make(chan V)
	_ = out2
	Apply(
		Of(
			V{1, "a"},
			V{2, "b"},
			V{1, "aa"},
			V{3, "c"},
			V{2, "bb"},
			V{3, "cc"},
			V{4, "d"},
		),
		Pipe3(
			GroupBy(func(v V) int {
				return v.Id
			}),
			// MergeAll[V](),
			MergeMap(func(g Observable[V]) Observable[[]V] {
				return Apply(
					g,
					Pipe1(
						ToSlice[V](),
					),
				)
			}),
			ToChan(out),
		),
	).Subscribe(NewSubscriberEmpty[[]V]())

	expect := [][]V{
		{{1, "a"}, {1, "aa"}},
		{{2, "b"}, {2, "bb"}},
		{{3, "c"}, {3, "cc"}},
		{{4, "d"}},
	}

	i := 0
	for o := range out {
		i++
		// log.Printf("%#v", o)
		require.Contains(t, expect, o)
	}
	require.Equal(t, len(expect), i)
}
