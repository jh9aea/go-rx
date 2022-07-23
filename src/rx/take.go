package rx

import (
	"sync"
)

type counter struct {
	count int
	m     sync.Mutex
}

func Take[A any](n int) Operator[A, A] {
	return func(in Observable[A]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			co := counter{count: 0}
			subs := newSubscriptions()
			subs.Add(in.Subscribe(
				NewSubscriberFn(
					func(v A) {
						co.m.Lock()
						defer co.m.Unlock()
						co.count++
						s.OnNext(v)

						if co.count >= n {
							s.OnComplete()
							subs.Unsubscribe()
						}
					},
					deletegateOnError(s),
					delegateOnComplete(s),
				),
			))
			return subs.Unsubscribe
		})
	}
}
