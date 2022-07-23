package rx

func Filter[T any](pred func(v T) (bool, error)) Operator[T, T] {
	return func(in Observable[T]) Observable[T] {
		return Create(func(s SubscriberType[T]) Unsubscribe {
			return in.Subscribe(NewSubscriberFn(
				func(v T) {
					ok, err := pred(v)
					if err != nil {
						s.OnError(err)
						return
					}
					if ok {
						s.OnNext(v)
					}
				},
				deletegateOnError(s),
				delegateOnComplete(s),
			))
		})
	}
}
