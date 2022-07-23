package rx

type (
	OnNext[T any] func(v T)
	OnError       func(e error)
	OnComplete    func()

	SubscriberType[T any] interface {
		OnNext(v T)
		OnError(e error)
		OnComplete()
	}

	subscriber[T any] struct {
		onN    OnNext[T]
		onErr  OnError
		onComp OnComplete
	}

	SubscriberArg[T any] struct {
		OnNext     OnNext[T]
		OnError    OnError
		OnComplete OnComplete
	}
)

func (s subscriber[T]) OnNext(v T) {
	s.onN(v)
}

func (s subscriber[T]) OnError(e error) {
	s.onErr(e)
}

func (s subscriber[T]) OnComplete() {
	s.onComp()
}

func NewSubscriber[T any](sa SubscriberArg[T]) SubscriberType[T] {
	OnNext := sa.OnNext
	OnError := sa.OnError
	OnComplete := sa.OnComplete

	if OnNext == nil {
		OnNext = func(v T) {}
	}

	if OnError == nil {
		OnError = func(e error) {}
	}

	if OnComplete == nil {
		OnComplete = func() {}
	}

	return subscriber[T]{
		onN:    OnNext,
		onErr:  OnError,
		onComp: OnComplete,
	}
}

func NewSubscriberNext[T any](onNext OnNext[T]) SubscriberType[T] {
	return NewSubscriber(SubscriberArg[T]{OnNext: onNext})
}

func NewSubscriberComplete[T any](onComplete OnComplete) SubscriberType[T] {
	return NewSubscriber(SubscriberArg[T]{OnComplete: onComplete})
}

func NewSubscriberFn[T any](onNext OnNext[T], onError OnError, onComp OnComplete) SubscriberType[T] {
	return NewSubscriber(SubscriberArg[T]{
		OnNext:     onNext,
		OnError:    onError,
		OnComplete: onComp,
	})
}

func NewSubscriberEmpty[T any](t ...any) SubscriberType[T] {
	return NewSubscriber(SubscriberArg[T]{})
}

func delegateOnNext[T any](s SubscriberType[T]) func(v T) {
	return func(v T) {
		s.OnNext(v)
	}
}

func deletegateOnError[T any](s SubscriberType[T]) func(e error) {
	return func(e error) {
		s.OnError(e)
	}
}

func delegateOnComplete[T any](s SubscriberType[T]) func() {
	return func() {
		s.OnComplete()
	}
}
