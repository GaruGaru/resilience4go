package resilience

type RFunc[T any] func() (T, error)

type RMiddlware[T any] interface {
	Execute(RFunc[T]) RFunc[T]
}

type Runner[T any] struct {
	middlewares []RMiddlware[T]
}

func New[T any](middlewares ...RMiddlware[T]) *Runner[T] {
	return &Runner[T]{middlewares: middlewares}
}

func (r *Runner[T]) Execute(fn RFunc[T]) (T, error) {
	if len(r.middlewares) == 0 {
		return fn()
	}
	var composed RFunc[T]
	for _, md := range r.middlewares {
		f := md.Execute
		if composed == nil {
			composed = f(fn)
			continue
		}
		composed = f(composed)
	}
	return composed()
}
