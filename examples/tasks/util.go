package main

func mapf[T any, R any](in []T, transform func(T) R) []R {
	out := make([]R, len(in))
	for i := range in {
		out[i] = transform(in[i])
	}
	return out
}

func filterf[T any](in []T, filter func(T) bool) []T {
	out := make([]T, 0)
	for i := range in {
		if filter(in[i]) {
			out = append(out, in[i])
		}
	}
	return out
}
