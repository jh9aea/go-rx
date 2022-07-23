package rx

import "sync"

func Buffer[A any](count int) Operator[A, []A] {
	return func(in Observable[A]) Observable[[]A] {
		return Create(func(s SubscriberType[[]A]) Unsubscribe {
			var buffer []A
			m := sync.Mutex{}
			return in.Subscribe(
				NewSubscriberFn(
					func(v A) {
						m.Lock()
						m.Unlock()
						buffer = append(buffer, v)
						if len(buffer) >= count {
							s.OnNext(buffer)
							buffer = nil
						}
					},
					func(e error) {
						s.OnError(e)
					},
					func() {
						m.Lock()
						if len(buffer) > 0 {
							s.OnNext(buffer)
						}
						m.Unlock()

						s.OnComplete()
					},
				),
			)
		})
	}
}
