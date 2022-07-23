package rx

import "sync"

// https://reactivex.io/documentation/operators/reduce.html
func Reduce[T, A any](fn func(acc A, v T) (A, error), initial A) Operator[T, A] {
	return func(in Observable[T]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			m := sync.Mutex{}
			cur := initial
			return in.Subscribe(NewSubscriberFn(
				func(v T) {
					var err error
					m.Lock()
					cur, err = fn(cur, v)
					m.Unlock()
					if err != nil {
						s.OnError(err)
						return
					}
				},
				deletegateOnError(s),
				func() {
					s.OnNext(cur)
					s.OnComplete()
				},
			))
		})
	}
}
