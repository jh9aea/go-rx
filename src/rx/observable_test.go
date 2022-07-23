package rx

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/require"
)

func TestCreateBlocking(t *testing.T) {
	i := 0
	Create(func(s SubscriberType[int]) Unsubscribe {
		s.OnNext(1)
		s.OnNext(2)
		s.OnComplete()
		return func() {}
	}).Subscribe(
		NewSubscriber[int](SubscriberArg[int]{
			OnNext: func(v int) {
				i++
				require.Equal(t, i, v)
			},
			OnComplete: func() {
				i++
			},
		}),
	)
	require.Equal(t, 3, i)
}

func TestCreateAsync(t *testing.T) {
	i := 0
	un := Apply(
		Create(func(s SubscriberType[int]) Unsubscribe {
			d := make(chan any, 1)
			go func() {
				select {
				case <-d:
					return
				case <-time.After(time.Nanosecond):
				}

				for i := 0; i < 10; i++ {
					select {
					case <-d:
						return
					default:
						s.OnNext(i)
					}
				}
				s.OnComplete()
			}()
			return func() {
				d <- true
			}
		}),
		Do(func(v int) { i = v }),
	).Subscribe(NewSubscriberEmpty[int]())
	un()
	require.Equal(t, 0, i)
}

func TestDeferred(t *testing.T) {
	wg := sync.WaitGroup{}
	source :=
		Defer(func() Observable[int] {
			return Of(rand.Int())
		})

	var v1, v2 int
	wg.Add(2)
	source.Subscribe(NewSubscriberNext(func(v int) {
		v1 = v
		wg.Done()
	}))
	source.Subscribe(NewSubscriberNext(func(v int) {
		v2 = v
		wg.Done()

	}))
	wg.Wait()
	require.NotEqual(t, v1, v2)

	source = Of(rand.Int())
	wg.Add(2)
	source.Subscribe(NewSubscriberNext(func(v int) { v1 = v; wg.Done() }))
	source.Subscribe(NewSubscriberNext(func(v int) { v2 = v; wg.Done() }))
	wg.Wait()
	require.Equal(t, v1, v2)

	wg.Add(2)
	Apply(
		source,
		Map(func(v int) (int, error) { return v - v, nil }),
	).
		Subscribe(NewSubscriberNext(func(v int) { v1 = v; wg.Done() }))
	source.Subscribe(NewSubscriberNext(func(v int) { v2 = v; wg.Done() }))
	wg.Wait()
	require.Equal(t, v1, 0)
	require.NotEqual(t, v1, v2)
}

var r int

func benchmarkObs(n int, b *testing.B) {
	for n := 0; n <= b.N; n++ {
		d, f := FinalizeCreateDone[int]()
		Apply(
			Range(0, n),
			Pipe2(
				f,
				Last[int](),
			),
		).Subscribe(NewSubscriberNext(func(v int) { r = v }))
		<-d
	}
}

func BenchmarkObservable1(b *testing.B) {
	benchmarkObs(1, b)
}

func BenchmarkObservable10(b *testing.B) {
	benchmarkObs(10, b)
}

func BenchmarkObservable100(b *testing.B) {
	benchmarkObs(100, b)
}

func BenchmarkObservable1000(b *testing.B) {
	benchmarkObs(1000, b)
}

func BenchmarkRxgo1000(b *testing.B) {
	for n := 0; n <= b.N; n++ {
		<-rxgo.Range(0, 1000).Last().Run()
	}
}
