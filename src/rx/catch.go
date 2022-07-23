package rx

import (
	"sync"
)

type cache_un struct {
	m sync.Mutex
	f func()
}

func (c *cache_un) Set(f func()) {
	c.m.Lock()
	defer c.m.Unlock()
	c.f = f
}

func (c *cache_un) Uns() {
	c.m.Lock()
	defer c.m.Unlock()
	c.f()
}

func Catch[A any](catch func(e error) Observable[A]) Operator[A, A] {
	return func(in Observable[A]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			cu := cache_un{}
			cu.Set(
				in.Subscribe(
					NewSubscriberFn(
						func(a A) { s.OnNext(a) },
						func(e error) {
							cu.Set(
								catch(e).Subscribe(NewSubscriberFn(
									func(b A) {
										s.OnNext(b)
									},
									func(e error) { s.OnError(e) },
									func() { s.OnComplete() },
								)),
							)
						},
						func() { s.OnComplete() },
					),
				),
			)
			return cu.Uns
		})
	}
}
