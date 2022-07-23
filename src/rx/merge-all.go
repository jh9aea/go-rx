package rx

import "sync"

// MergeAll
//	Apply(
//	 	Of(
//	 		Of(1),
//	 		Of(2),
//	 		Range(3, 4),
//	 		Apply(
//	 			Interval(10*time.Millisecond),
//	 			Take[int](2),
//	 		),
//	 	),
//	 	Pipe3(
//	 		MergeAll[int](),
//	 		ToSlice[int](),
//	 		ToChan(out),
//	 	),
//	 ).Subscribe(NewSubscriberEmpty[[]int]())
func MergeAll[A any]() Operator[Observable[A], A] {
	return func(in Observable[Observable[A]]) Observable[A] {
		return Create(func(s SubscriberType[A]) Unsubscribe {
			subs := newSubscriptions()
			wg := sync.WaitGroup{}
			wg.Add(1)

			subs.Add(in.Subscribe(NewSubscriberFn(
				func(v Observable[A]) {
					wg.Add(1)
					subs.Add(
						v.Subscribe(NewSubscriberFn(
							func(v A) {
								s.OnNext(v)
							},
							func(e error) {
								s.OnError(e)
							},
							func() {
								wg.Done()
								// mb need to remove from subs?
							},
						)),
					)
				},
				func(e error) {
					s.OnError(e)
				},
				func() {
					wg.Done()
				},
			)))

			go func() {
				wg.Wait()
				s.OnComplete()
			}()

			return subs.Unsubscribe
		})
	}
}
