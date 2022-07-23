package rx

import (
	"sync/atomic"
)

type Observable[T any] interface {
	Subscribe(SubscriberType[T]) Unsubscribe
}

type observable[T any] struct {
	gen func(SubscriberType[T]) Unsubscribe
}

func (o *observable[T]) Subscribe(subscriber SubscriberType[T]) Unsubscribe {
	subs := newSubscriptions()
	var errored int64 = 0
	subs.Add(o.gen(NewSubscriberFn(
		func(v T) {
			if subs.IsUnsubscribed() {
				return
			}
			subscriber.OnNext(v)
		},
		func(e error) {
			if subs.IsUnsubscribed() {
				return
			}
			subscriber.OnError(e)
			atomic.AddInt64(&errored, 1)
			subs.Unsubscribe()
		},
		func() {
			if subs.IsUnsubscribed() {
				return
			}
			if atomic.LoadInt64(&errored) == 0 {
				subscriber.OnComplete()
				subs.Unsubscribe()
			}
		},
	)))
	return subs.Unsubscribe
}

// Creates Observable
func Create[T any](fn func(SubscriberType[T]) Unsubscribe) Observable[T] {
	return &observable[T]{gen: fn}
}

func Defer[T any](fn func() Observable[T]) Observable[T] {
	return &observable[T]{gen: func(observer SubscriberType[T]) Unsubscribe {
		return fn().Subscribe(observer)
	}}
}

func Error[T any](e error) Observable[T] {
	return Create(func(s SubscriberType[T]) Unsubscribe {
		go s.OnError(e)
		return func() {}
	})
}

func Empty[T any]() Observable[T] {
	return Create(func(s SubscriberType[T]) Unsubscribe {
		go s.OnComplete()
		return func() {}
	})
}
