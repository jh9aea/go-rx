package rx

// draft
//  - should block?
//  - should close on unsubscribe?
func ToChan[A any](ch chan<- A) Operator[A, A] {
	return func(in Observable[A]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			// we need proxy for channel
			// because channel may be unbuffered
			// and it will cause to block
			que := make(chan A, 1)

			go func() {
				for {
					select {
					case v, ok := <-que:
						if !ok {
							close(ch)
							return
						}
						ch <- v
					}
				}
			}()

			return in.Subscribe(NewSubscriberFn(
				func(v A) {
					que <- v
					s.OnNext(v)
				},
				func(e error) {
					close(que)
					s.OnError(e)
				},
				func() {
					close(que)
					s.OnComplete()
				},
			))
		})
	}
}
