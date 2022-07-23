package rx

import "sync"

// https://reactivex.io/documentation/operators/concat.html
func Concat[T any](o ...Observable[T]) Observable[T] {
	return Create(func(s SubscriberType[T]) Unsubscribe {
		subs := newSubscriptions()
		wg := sync.WaitGroup{}

		go func() {
			for _, obs := range o {
				wg.Add(1)
				subs.Add(
					obs.Subscribe(NewSubscriberFn(
						delegateOnNext(s),
						func(e error) {
							s.OnError(e)
							wg.Done()
						},
						func() {
							wg.Done()
						}),
					),
				)
				wg.Wait()
			}

			s.OnComplete()
			subs.Unsubscribe()
		}()

		return subs.Unsubscribe
	})
}
