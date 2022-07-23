package rx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMergeConcatfunc(t *testing.T) {
	complete := 0
	r := []int{}
	done, fin := FinalizeCreateDone[[]int]()

	Apply(
		Merge(
			Interval(time.Millisecond*8),
			Apply(Interval(time.Millisecond*46), Map(func(v int) (int, error) { return 200 + v, nil })),
		),
		Pipe4(
			TakeUntil[int](Interval(time.Millisecond*47)),
			ToSlice[int](),
			Do(func(v []int) {
				r = v
			}),
			fin,
		),
	).Subscribe(
		NewSubscriberFn(
			func(v []int) {},
			func(e error) { t.Error("should not fire") },
			func() { complete++ },
		),
	)
	<-done

	require.Equal(t, []int{0, 1, 2, 3, 4, 200}, r)
	require.Equal(t, 1, complete)
}

func TestMergeParallelfunc(t *testing.T) {
	complete := 0
	r := []int{}
	done, fin := FinalizeCreateDone[[]int]()

	Apply(
		Merge(
			Interval(time.Millisecond*11),
			Apply(
				Interval(time.Millisecond*20),
				Pipe1(
					Map(func(v int) (int, error) { return 200 + v, nil }),
				),
			),
		),
		Pipe4(
			TakeUntil[int](Interval(time.Millisecond*30)),
			ToSlice[int](),
			Do(func(v []int) { r = v }),
			fin,
		),
	).Subscribe(NewSubscriberFn(func(v []int) {},
		func(e error) { t.Error("should not fire") },
		func() { complete++ }))
	<-done

	require.Equal(t, []int{0, 200, 1}, r)
	require.Equal(t, 1, complete)
}

func TestMergeComplete(t *testing.T) {
	complete := 0
	r := []int{}
	done, fin := FinalizeCreateDone[[]int]()

	Apply(
		Merge(
			Of(1, 2),
			Of(20, 30),
		),
		Pipe3(
			ToSlice[int](),
			Do(func(v []int) { r = v }),
			fin,
		),
	).Subscribe(NewSubscriberFn(func(v []int) {}, func(e error) { t.Error("should not fire") }, func() { complete++ }))
	<-done

	require.Equal(t, 4, len(r))
	require.Equal(t, 1, complete)
}
