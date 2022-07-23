package rx

import "sync"

// TODO: mb should be an interface Subscription{Unsubscribe(); IsSubscribed() bool}
//
type Unsubscribe func()

type subscriptions struct {
	subs []Unsubscribe
	l    sync.Mutex
	uns  bool
}

func newSubscriptions() *subscriptions {
	return &subscriptions{subs: []Unsubscribe{}, uns: false}
}

func (s *subscriptions) Add(fs ...Unsubscribe) {
	s.l.Lock()
	defer s.l.Unlock()
	for _, un := range fs {
		s.subs = append(s.subs, un)
	}
}

func (s *subscriptions) Unsubscribe() {
	s.l.Lock()
	defer s.l.Unlock()

	s.uns = true
	for _, v := range s.subs {
		v()
	}
	s.subs = []Unsubscribe{}
}

func (s *subscriptions) IsUnsubscribed() bool {
	s.l.Lock()
	defer s.l.Unlock()

	return s.uns == true
}
