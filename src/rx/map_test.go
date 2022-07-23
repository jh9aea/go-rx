package rx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	done := make(chan bool, 1)
	r := []int{}
	Apply(
		Of(1, 2, 3),
		Pipe2(
			Map(func(a int) (int, error) { return a * a, nil }),
			ToSlice[int](),
		),
	).Subscribe(NewSubscriberFn(func(v []int) {
		r = v
	}, func(e error) {}, func() { close(done) }))
	<-done
	require.Equal(t, []int{1, 4, 9}, r)
}

func TestMapError(t *testing.T) {
	done := make(chan bool, 1)
	r := []int{}
	var err error
	Map(func(a int) (int, error) {
		if a == 2 {
			return 0, errors.New("")
		}
		return a * a, nil
	})(Of(1, 2, 3)).Subscribe(NewSubscriberFn(
		func(v int) {
			r = append(r, v)
		},
		func(e error) {
			err = e
			close(done)
		},
		func() {
			close(done)
		}),
	)
	<-done
	require.Equal(t, []int{1}, r)
	require.Error(t, err)
}
