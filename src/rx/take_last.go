package rx

import "sync"

func TakeLast[T any](n int) Operator[T, []T] {
	return func(in Observable[T]) Observable[[]T] {
		return Create(func(s SubscriberType[[]T]) Unsubscribe {
			lastN := make([]T, n)
			m := sync.Mutex{}
			return in.Subscribe(NewSubscriberFn(
				func(v T) {
					m.Lock()
					defer m.Unlock()
					lastN = append(lastN[1:], v)
				},
				func(e error) { s.OnError(e) },
				func() {
					s.OnNext(lastN)
					s.OnComplete()
				},
			))
		})
	}
}
