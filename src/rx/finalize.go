package rx

func Finalize[A any](fn func()) Operator[A, A] {
	return func(in Observable[A]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			return in.Subscribe(NewSubscriberFn(
				delegateOnNext(s),
				func(e error) {
					s.OnError(e)
					fn()
				},
				func() {
					s.OnComplete()
					fn()
				},
			))
		})
	}
}

func FinalizeCreateDone[A any]() (<-chan any, Operator[A, A]) {
	done := make(chan any, 1)
	return done, Finalize[A](func() { done <- struct{}{} })
}
