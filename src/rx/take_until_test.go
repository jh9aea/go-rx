package rx

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTakeUntil(t *testing.T) {
	done := make(chan any)
	r := 0
	Apply(
		Apply(
			Interval(time.Millisecond*100),
			Do(func(v int) {
				// log.Printf("%#v", v)
				r++
			}),
		),
		Pipe2(
			TakeUntil[int](
				Interval(time.Millisecond*350),
			),
			Finalize[int](func() { done <- 1 }),
		),
	).Subscribe(NewSubscriberEmpty[int]())
	<-done
	require.Equal(t, 3, r)
}

func TestTakeUntilCom(t *testing.T) {
	d := make(chan any)
	m := sync.Mutex{}
	r := []int{}
	Apply(
		Range(0, 1),
		TakeUntil[int](
			Interval(time.Second*100),
		),
	).Subscribe(NewSubscriberFn(
		func(v int) {
			m.Lock()
			defer m.Unlock()
			r = append(r, v)
		},
		func(e error) { t.Error("should not") },
		func() {
			d <- true
		},
	))
	<-d
	require.Equal(t, []int{0, 1}, r)
}

func TestTakeUntilCtx(t *testing.T) {
	d := make(chan any)
	var r []int
	c := 0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Apply(
		Interval(time.Millisecond*10),
		Pipe3(
			Do(func(v int) {
				if v == 4 {
					cancel()
				}
			}),
			TakeUntilCtx[int](ctx),
			ToSlice[int](),
		),
	).Subscribe(NewSubscriberFn(
		func(v []int) { r = v },
		func(e error) {},
		func() {
			c++
			d <- true
		},
	))
	<-d
	require.Equal(t, 1, c)
	require.Equal(t, []int{0, 1, 2, 3, 4}, r)
}

func TestTakeUntilCtxComplete(t *testing.T) {
	d := make(chan any)
	var r []int
	c := 0

	Apply(
		Range(0, 4),
		Pipe2(
			TakeUntilCtx[int](context.Background()),
			ToSlice[int](),
		),
	).Subscribe(NewSubscriberFn(
		func(v []int) { r = v },
		func(e error) {},
		func() {
			c++
			d <- true
		},
	))
	<-d
	require.Equal(t, 1, c)
	require.Equal(t, []int{0, 1, 2, 3, 4}, r)
}

func TestTakeUntilCtxUns(t *testing.T) {
	d := make(chan any)
	m := sync.Mutex{}
	var r []int
	c := 0

	uns := Apply(
		Interval(time.Millisecond*10),
		Pipe2(
			TakeUntilCtx[int](context.Background()),
			Do(func(v int) { m.Lock(); defer m.Unlock(); r = append(r, v) }),
		),
	).Subscribe(NewSubscriberFn(
		func(v int) {},
		func(e error) {},
		func() { m.Lock(); defer m.Unlock(); c++ },
	))

	go func() {
		time.Sleep(time.Millisecond * 55)
		uns()
		d <- true
	}()

	<-d
	m.Lock()
	defer m.Unlock()

	require.Equal(t, 0, c)
	require.Equal(t, []int{0, 1, 2, 3, 4}, r)
}
