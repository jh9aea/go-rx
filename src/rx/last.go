package rx

import "sync"

func Last[T any]() Operator[T, T] {
	return func(in Observable[T]) Observable[T] {
		return Create(func(s SubscriberType[T]) Unsubscribe {
			m := sync.Mutex{}
			var last T
			return in.Subscribe(NewSubscriberFn(
				func(v T) {
					m.Lock()
					defer m.Unlock()
					last = v
				},
				func(e error) { s.OnError(e) },
				func() {
					s.OnNext(last)
					s.OnComplete()
				},
			))
		})
	}
}
