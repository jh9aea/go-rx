package rx

type Operator[A, B any] func(Observable[A]) Observable[B]
