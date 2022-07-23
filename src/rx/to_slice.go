package rx

import "sync"

func ToSlice[T any]() Operator[T, []T] {
	return func(in Observable[T]) Observable[[]T] {
		return Create(func(s SubscriberType[[]T]) Unsubscribe {
			r := []T{}
			m := sync.Mutex{}
			return in.Subscribe(NewSubscriberFn(
				func(v T) {
					m.Lock()
					defer m.Unlock()
					r = append(r, v)
				},
				func(err error) {
					s.OnError(err)
				},
				func() {
					m.Lock()
					defer m.Unlock()
					s.OnNext(r)
					s.OnComplete()
				},
			))
		})
	}
}
