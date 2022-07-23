package rx

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShare(t *testing.T) {
	n := 4
	multi := Apply(
		Create(func(s SubscriberType[int]) Unsubscribe {
			go func() {
				time.Sleep(1 * time.Millisecond)
				for i := 0; i < n; i++ {
					s.OnNext(rand.Int())
				}
				s.OnComplete()
			}()
			return func() {}
		}),
		Share[int](),
	)

	cout1, cout2 := make(chan int, n), make(chan int, n)

	Apply(multi, ToChan(cout1)).
		Subscribe(NewSubscriberEmpty[int]())
	Apply(multi, ToChan(cout2)).
		Subscribe(NewSubscriberEmpty[int]())

	out1, out2 := []int{}, []int{}
	for v := range cout1 {
		out1 = append(out1, v)
	}
	for v := range cout2 {
		out2 = append(out2, v)
	}

	require.Equal(t, n, len(out1))
	require.Equal(t, n, len(out2))

	for i := range out1 {
		require.Equal(t, out1[i], out2[i])
	}
}
