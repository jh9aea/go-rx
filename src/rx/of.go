package rx

func Of[T any](v ...T) Observable[T] {
	return Create(func(s SubscriberType[T]) Unsubscribe {
		done := make(chan any, 1)
		go func() {
			func() {
				for _, v := range v {
					select {
					case <-done:
						return
					default:
						s.OnNext(v)
					}
				}
			}()
			s.OnComplete()
		}()
		return func() {
			done <- struct{}{}
		}
	})
}

func OfSlice[T any](v []T) Observable[T] {
	return Of(v...)
}

func Range(from, to int) Observable[int] {
	return Create(func(s SubscriberType[int]) Unsubscribe {
		done := make(chan bool, 1)

		go func() {
			for i := from; i <= to; i++ {
				select {
				case <-done:
					return
				default:
					s.OnNext(i)
				}
			}
			s.OnComplete()
		}()

		return func() {
			done <- true
		}
	})
}
