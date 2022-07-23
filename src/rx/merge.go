package rx

import (
	"sync"
)

func Merge[T any](obs ...Observable[T]) Observable[T] {
	return Create(func(s SubscriberType[T]) Unsubscribe {
		wg := sync.WaitGroup{}
		wg.Add(len(obs))
		subs := newSubscriptions()

		for _, o := range obs {
			subs.Add(
				o.Subscribe(NewSubscriberFn(
					func(v T) {
						s.OnNext(v)
					},
					func(e error) {
						s.OnError(e)
						wg.Done()
					},
					func() {
						wg.Done()
					},
				)),
			)
		}

		go func() {
			wg.Wait()
			s.OnComplete()
			subs.Unsubscribe()
		}()

		return subs.Unsubscribe
	})
}
