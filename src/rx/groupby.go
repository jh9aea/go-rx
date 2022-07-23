package rx

import "sync"

// https://reactivex.io/documentation/operators/groupby.html
func GroupBy[A any, B comparable](groupper func(v A) B /*, error ? */) Operator[A, Observable[A]] {
	return func(in Observable[A]) Observable[Observable[A]] {
		return Create(func(s SubscriberType[Observable[A]]) Unsubscribe {
			groups := map[B]*Subject[A]{}
			m := sync.Mutex{}
			return in.Subscribe(NewSubscriberFn(
				func(v A) {
					groupId := groupper(v)

					m.Lock()
					gr, ok := groups[groupId]
					if !ok {
						gr = NewSubject[A]()
						groups[groupId] = gr
						s.OnNext(gr.Observable())
					}
					m.Unlock()

					gr.OnNext(v)
				},
				func(e error) {
					m.Lock()
					defer m.Unlock()
					for _, gr := range groups {
						gr.OnError(e)
					}
				},
				func() {
					m.Lock()
					defer m.Unlock()
					for _, gr := range groups {
						gr.OnComplete()
					}
					s.OnComplete()
				},
			))
		})
	}
}
