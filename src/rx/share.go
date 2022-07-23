package rx

import "sync"

// https://reactivex.io/documentation/operators/refcount.html
func Share[A any]() Operator[A, A] {
	return func(in Observable[A]) Observable[A] {
		sub := NewSubject[A]()
		o := sync.Once{}
		return Defer(func() Observable[A] {
			o.Do(func() {
				in.Subscribe(sub)
			})
			return sub
		})
	}
}
