package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubject(t *testing.T) {
	out1 := make(chan []int, 4)
	out2 := make(chan int, 4)
	out3 := make(chan int, 4)

	s := NewSubject[int]()

	o1 := Apply[int](
		s,
		Pipe4(
			Take[int](2),
			Map(func(v int) (int, error) { return v * v, nil }),
			ToSlice[int](),
			ToChan[[]int](out1),
		),
	).Subscribe(NewSubscriberEmpty[[]int]())

	o2 := Apply[int](
		s,
		Pipe3(
			Filter(func(v int) (bool, error) {
				return v > 2, nil
			}),
			Scan(func(acc, v int) (int, error) {
				return acc + v, nil
			}, 0),
			ToChan[int](out2),
		),
	).Subscribe(NewSubscriberEmpty[int]())

	_, _ = o1, o2

	go func() {
		s.OnNext(1)
		s.OnNext(2)

		o3 := Apply[int](s, ToChan(out3)).
			Subscribe(NewSubscriberEmpty[int]())

		s.OnNext(3)
		o3()

		s.OnNext(4)
		s.OnComplete()
	}()

	v1 := <-out1
	require.Equal(t, []int{1, 4}, v1)

	i := 0
	r := []int{3, 7}
	for v2 := range out2 {
		require.Equal(t, r[i], v2)
		i++
	}

	// should out3 be closed on unsubscribe toChan
	v3 := <-out3
	require.Equal(t, 3, v3)
}

func TestSubjectComplete(t *testing.T) {
	s := NewSubject[any]()
	s.OnComplete()
	c := 0
	d := make(chan any)

	Apply[any](
		s,
		Finalize[any](func() { d <- true }),
	).Subscribe(NewSubscriberFn[any](
		func(v any) { t.Error("should not OnNext") },
		func(e error) { t.Error("should not OnError") },
		func() { c++ },
	))

	<-d

	require.Equal(t, 1, c)
}
