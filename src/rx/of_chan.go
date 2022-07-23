package rx

// Create cold observable from channel
// It will not start <- chan until .Subscribe()
func OfChan[T any](ch <-chan T) Observable[T] {
	return Create(func(s SubscriberType[T]) Unsubscribe {
		done := make(chan struct{}, 1)
		go func() {
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						s.OnComplete()
						return
					}
					s.OnNext(v)
				case <-done:
					return
				}
			}
		}()
		return func() {
			done <- struct{}{}
		}
	})
}
