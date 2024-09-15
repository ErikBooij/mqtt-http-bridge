package utilities

func FilterSlice[T any](s []T, f func(T) bool) []T {
	out := make([]T, 0)

	for _, v := range s {
		if f(v) {
			out = append(out, v)
		}
	}

	return out
}

func MapSlice[In any, Out any](s []In, f func(In) Out) []Out {
	out := make([]Out, len(s))

	for i, v := range s {
		out[i] = f(v)
	}

	return out
}
