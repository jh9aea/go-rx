package rx

import (
	"context"
)

// https://reactivex.io/documentation/operators/takeuntil.html
func TakeUntil[A, B any](un Observable[B]) Operator[A, A] {
	return func(in Observable[A]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			subs := newSubscriptions()
			subs.Add(
				in.Subscribe(s),
				un.Subscribe(NewSubscriberFn(
					func(B) {
						s.OnComplete()
						subs.Unsubscribe()
					},
					func(e error) { s.OnError(e) },
					func() {},
				)),
			)
			return subs.Unsubscribe
		})
	}
}

func TakeUntilChan[A, C any](c <-chan C) Operator[A, A] {
	return TakeUntil[A](Create(func(s SubscriberType[any]) Unsubscribe {
		d := make(chan any)
		go func() {
			select {
			case <-c:
				s.OnNext(struct{}{})
			case <-d:
			}
		}()
		return func() {
			d <- struct{}{}
		}
	}))
}

func TakeUntilCtx[A any](ctx context.Context) Operator[A, A] {
	return TakeUntilChan[A](ctx.Done())
}
