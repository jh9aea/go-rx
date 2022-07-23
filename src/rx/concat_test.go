package rx

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConcat(t *testing.T) {
	var c int64 = 0
	out := make(chan []int)

	Apply(
		Concat(
			Apply(Interval(time.Millisecond*5), Pipe2(Map(func(v int) (int, error) { return v + 100, nil }), Take[int](2))),
			Apply(Interval(time.Millisecond*1), Pipe2(Map(func(v int) (int, error) { return v + 200, nil }), Take[int](2))),
		),
		Pipe3(
			Map(func(v int) (int, error) { return v + v, nil }),
			ToSlice[int](),
			ToChan(out),
		),
	).Subscribe(NewSubscriberFn(func(v []int) {}, func(e error) { t.Error("should not be") }, func() { atomic.AddInt64(&c, 1) }))

	r := <-out
	require.Equal(t, []int{200, 202, 400, 402}, r)
	require.Equal(t, int64(1), atomic.LoadInt64(&c))
}

func TestConcatError(t *testing.T) {
	c := 0
	out := make(chan []int)
	err := make(chan error)

	Apply(
		Concat(
			Apply(Interval(time.Millisecond*5), Pipe2(Map(func(v int) (int, error) { return v + 100, nil }), Take[int](2))),
			Apply(
				Interval(time.Millisecond*1),
				Pipe2(Map(func(v int) (int, error) { return v + 200, errors.New("EEE") }), Take[int](2)),
			),
		),
		Pipe3(
			Map(func(v int) (int, error) { return v + v, nil }),
			ToSlice[int](),
			ToChan(out),
		),
	).Subscribe(NewSubscriberFn(
		func(v []int) { t.Error("should not") },
		func(e error) {
			err <- e
		},
		func() { c++ },
	))
	<-out

	require.Equal(t, 0, c)
	require.Equal(t, "EEE", (<-err).Error())
}

func TestConcatUnsub(t *testing.T) {
	d := make(chan any, 1)
	m := sync.Mutex{}
	r := []int{}

	uns := Concat(
		Apply(Interval(time.Microsecond*10), Take[int](2)),
		Apply(Interval(time.Microsecond*20), Pipe[int]()),
	).Subscribe(NewSubscriberFn(
		func(v int) {
			m.Lock()
			defer m.Unlock()
			r = append(r, v)
		},
		func(e error) { t.Error("should not error") },
		func() { t.Error("should not complete") },
	))

	go func() {
		time.Sleep(100 * time.Millisecond)
		uns()
		d <- 1
	}()

	<-d
}
