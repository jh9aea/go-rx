package rx

import "sync"

// https://reactivex.io/documentation/operators/skipuntil.html
func SkipUntil[T any, O any](other Observable[O]) Operator[T, T] {
	return func(in Observable[T]) Observable[T] {
		return Create(func(s SubscriberType[T]) Unsubscribe {
			m := sync.Mutex{}
			run := false
			var err error
			subs := newSubscriptions()
			subs.Add(
				in.Subscribe(
					NewSubscriberFn[T](
						func(v T) {
							m.Lock()
							defer m.Unlock()
							if run {
								s.OnNext(v)
							}
						},
						func(e error) {
							m.Lock()
							defer m.Unlock()
							err = e
							if run {
								s.OnError(err)
							}
						},
						func() {
							s.OnComplete()
						},
					),
				),
				other.Subscribe(
					NewSubscriberNext[O](func(v O) {
						m.Lock()
						defer m.Unlock()
						run = true
						if err != nil {
							s.OnError(err)
						}
					}),
				),
			)
			return subs.Unsubscribe
		})
	}
}
