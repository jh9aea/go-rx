package rx

import (
	"context"
	"sync"
)

// Draft
//
// will not respect the order
//
// 	go Apply(
// 		Range(0, 5),
// 		Pipe2(
// 			PoolInfinite(
// 				Map(func(v int) (int, error) {
// 					time.Sleep(time.Duration(250-v*25) * time.Millisecond)
// 					return v, nil
// 				}),
// 			),
// 			ToSlice[int](),
// 	 	),
// 	)
func PoolInfinite[A, B any](o Operator[A, B]) Operator[A, B] {
	return func(in Observable[A]) Observable[B] {
		return Create(func(s SubscriberType[B]) Unsubscribe {
			wg := sync.WaitGroup{}
			// for OnComplete
			wg.Add(1)
			ctx, cancel := context.WithCancel(context.Background())
			un := Apply(
				in,
				TakeUntilCtx[A](ctx),
			).
				Subscribe(NewSubscriberFn(
					func(v A) {
						// this is dangerous
						wg.Add(1)
						go Apply(
							Of(v),
							Pipe2(
								TakeUntilCtx[A](ctx),
								o,
							),
						).
							Subscribe(NewSubscriberFn(
								func(v B) {
									s.OnNext(v)
								},
								func(e error) {
									s.OnError(e)
								},
								func() {
									wg.Done()
								},
							))
					},
					func(e error) {
						s.OnError(e)
					},
					func() {
						wg.Done()
					},
				))
			go func() {
				wg.Wait()
				s.OnComplete()
			}()
			return func() {
				un()
				cancel()
			}
		})
	}
}

// Draft
func Pool[A, B any](num int, o Operator[A, B]) Operator[A, B] {
	return func(in Observable[A]) Observable[B] {
		return Create(func(s SubscriberType[B]) Unsubscribe {
			wg := sync.WaitGroup{}
			// for OnComplete
			wg.Add(1)
			ctx, cancel := context.
				WithCancel(context.Background())
			pool := make(chan struct{}, num)
			un := Apply(
				in,
				TakeUntilCtx[A](ctx),
			).
				Subscribe(NewSubscriberFn(
					func(v A) {
						// this is dangerous
						wg.Add(1)
						pool <- struct{}{}
						go Apply(
							Of(v),
							Pipe2(
								TakeUntilCtx[A](ctx),
								o,
							),
						).
							Subscribe(NewSubscriberFn(
								func(v B) {
									s.OnNext(v)
								},
								func(e error) {
									s.OnError(e)
								},
								func() {
									<-pool
									wg.Done()
								},
							))
					},
					func(e error) {
						s.OnError(e)
					},
					func() {
						wg.Done()
					},
				))
			go func() {
				wg.Wait()
				s.OnComplete()
			}()
			return func() {
				un()
				cancel()
			}
		})
	}
}
