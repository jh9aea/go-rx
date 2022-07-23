package rx

// https://reactivex.io/documentation/operators/flatmap.html
func MergeMap[A, B any](mapper func(v A) Observable[B]) Operator[A, B] {
	return func(in Observable[A]) Observable[B] {
		return Apply(
			in,
			Pipe2(
				Map(func(v A) (Observable[B], error) {
					return mapper(v), nil
				}),
				MergeAll[B](),
			),
		)
	}
}
