package rx

func Map[A, B any](fn func(a A) (B, error)) Operator[A, B] {
	return func(in Observable[A]) Observable[B] {
		return Create(func(s SubscriberType[B]) Unsubscribe {
			return in.Subscribe(NewSubscriberFn(
				func(a A) {
					v, err := fn(a)
					if err != nil {
						s.OnError(err)
						return
					}
					s.OnNext(v)
				},
				deletegateOnError(s),
				delegateOnComplete(s),
			))
		})
	}
}
