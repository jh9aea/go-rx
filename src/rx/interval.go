package rx

import (
	"time"
)

func Interval(t time.Duration) Observable[int] {
	return Create(func(s SubscriberType[int]) Unsubscribe {
		i := 0
		done := make(chan struct{}, 1)

		go func() {
			func() {
				for {
					select {
					case <-done:
						return
					case <-time.After(t):
						s.OnNext(i)
						i++
					}
				}
			}()
		}()

		return func() {
			done <- struct{}{}
		}
	})
}
