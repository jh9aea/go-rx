package rx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCatch(t *testing.T) {
	var r []int
	done := make(chan any)
	Apply(
		Of(1, 2, 3),
		Pipe2(
			Map(func(v int) (int, error) {
				if v >= 2 {
					return v, errors.New(">=2")
				}
				return v, nil
			}),
			Catch(func(e error) Observable[int] {
				return Of(100, 200)
			}),
		),
	).Subscribe(
		NewSubscriberFn(
			func(v int) {
				r = append(r, v)
			},
			func(e error) {
			},
			func() {
				done <- true
			},
		),
	)
	<-done
	require.Equal(t, []int{1, 100, 200}, r)
}
