package deep

func Append[T any](slice []T, items ...T) []T {
	result := make([]T, 0, len(slice)+len(items))

	return append(append(result, slice...), items...)
}
