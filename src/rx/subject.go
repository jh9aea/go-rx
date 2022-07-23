package rx

import "sync"

type Subject[T any] struct {
	subs map[int]SubscriberType[T]
	m    sync.RWMutex

	i int

	err          error
	compeleted   bool
	unsubscribed bool
}

func NewSubject[T any]() *Subject[T] {
	return &Subject[T]{subs: map[int]SubscriberType[T]{}}
}

func (s *Subject[T]) Stopped() bool {
	return s.compeleted || s.unsubscribed
}

func (s *Subject[T]) curSubs() []SubscriberType[T] {
	var r []SubscriberType[T]
	for _, s := range s.subs {
		r = append(r, s)
	}
	return r
}

func (s *Subject[T]) OnNext(v T) {
	s.m.RLock()
	subs := s.curSubs()
	s.m.RUnlock()
	for _, sub := range subs {
		sub.OnNext(v)
	}
}

func (s *Subject[T]) OnError(e error) {
	s.m.Lock()
	s.err = e
	subs := s.curSubs()
	s.m.Unlock()
	for _, sub := range subs {
		sub.OnError(e)
	}
}

func (s *Subject[T]) OnComplete() {
	s.m.Lock()
	defer s.m.Unlock()
	for _, sub := range s.subs {
		go sub.OnComplete()
	}
	s.compeleted = true
}

func (s *Subject[T]) Subscribe(st SubscriberType[T]) Unsubscribe {
	s.m.Lock()
	defer s.m.Unlock()

	s.panicIfUnsubscribed()

	s.i++
	i := s.i
	s.subs[i] = st

	if s.err != nil {
		go st.OnError(s.err)
	} else if s.compeleted {
		go st.OnComplete()
	}

	return func() {
		s.m.Lock()
		defer s.m.Unlock()
		delete(s.subs, i)
	}
}

func (s *Subject[T]) Unsubscribe() {
	s.m.Lock()
	defer s.m.Unlock()
	s.panicIfUnsubscribed()
	s.unsubscribed = true
}

func (s *Subject[T]) panicIfUnsubscribed() {
	if s.unsubscribed {
		panic("unsubscribed")
	}
}

func (s *Subject[T]) Observable() Observable[T] { return s }
