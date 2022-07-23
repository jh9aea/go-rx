package rx

// Applies operator to the observable
// Apply(Of(1,2), Take(1)).Subscribe()
func Apply[A, B any](o Observable[A], op Operator[A, B]) Observable[B] {
	return op(o)
}

func Pipe[A any]() Operator[A, A] { return func(o Observable[A]) Observable[A] { return o } }

func Pipe1[A, B any](op Operator[A, B]) Operator[A, B] {
	return func(o Observable[A]) Observable[B] {
		return op(o)
	}
}

func Pipe2[A, B, C any](op1 Operator[A, B], op2 Operator[B, C]) Operator[A, C] {
	return func(o Observable[A]) Observable[C] {
		return op2(op1(o))
	}
}

func Pipe3[A, B, C, D any](op1 Operator[A, B], op2 Operator[B, C], op3 Operator[C, D]) Operator[A, D] {
	return func(o Observable[A]) Observable[D] {
		return op3(op2(op1(o)))
	}
}

func Pipe4[A, B, C, D, E any](op1 Operator[A, B], op2 Operator[B, C], op3 Operator[C, D], op4 Operator[D, E]) Operator[A, E] {
	return func(o Observable[A]) Observable[E] {
		return op4(op3(op2(op1(o))))
	}
}

func PipeAny(opers ...Operator[any, any]) Operator[any, any] {
	return func(in Observable[any]) Observable[any] {
		r := in
		for _, o := range opers {
			r = o(r)
		}
		return r
	}
}
