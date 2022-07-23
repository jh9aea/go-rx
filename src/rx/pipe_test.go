package rx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	Apply(
		Of(1, 2, 3),
		ToSlice[int](),
	).Subscribe(
		NewSubscriberNext(func(v []int) {
			require.Equal(t, []int{1, 2, 3}, v)
		}),
	)
}
