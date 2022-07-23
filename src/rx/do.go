package rx

import "sync"

// https://reactivex.io/documentation/operators/do.html
func Do[A any](fn func(A)) Operator[A, A] {
	return func(in Observable[A]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			return in.Subscribe(
				NewSubscriberFn(
					func(v A) {
						fn(v)
						s.OnNext(v)
					},
					deletegateOnError(s),
					delegateOnComplete(s),
				),
			)
		})
	}
}

// https://reactivex.io/documentation/operators/do.html
func DoOnFirst[A any](fn func(A)) Operator[A, A] {
	o := sync.Once{}
	return Do(func(v A) {
		o.Do(func() { fn(v) })
	})
}
