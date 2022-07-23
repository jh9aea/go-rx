package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOfChan(t *testing.T) {
	ch := make(chan any)
	done := make(chan bool)

	r := []int{2, 4}
	i := 0

	Apply(
		OfChan(ch),
		PipeAny(
			Take[any](2),
			Map(func(v any) (any, error) {
				vv := v.(int)
				return vv * 2, nil
			}),
		),
	).Subscribe(NewSubscriberFn(
		func(v any) {
			require.Equal(t, r[i], v)
			i++
		},
		func(e error) {},
		func() {
			i = 200
			done <- true
		},
	))

	go func() {
		ch <- 1
		ch <- 2
	}()

	<-done
	require.Equal(t, 200, i)
}

func TestChanClose(t *testing.T) {
	ch := make(chan any)
	done := make(chan bool)

	r := []int{2, 4}
	i := 0

	Apply(
		OfChan(ch),
		PipeAny(
			Map(func(v any) (any, error) {
				vv := v.(int)
				return vv * 2, nil
			}),
		),
	).Subscribe(NewSubscriberFn(
		func(v any) {
			require.Equal(t, r[i], v)
			i++
		},
		func(e error) {},
		func() {
			i = 200
			done <- true
		},
	))

	go func() {
		ch <- 1
		ch <- 2
		close(ch)
	}()

	<-done
	require.Equal(t, 200, i)
}
