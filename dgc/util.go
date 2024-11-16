package dgc

func removeElement[T comparable](slice []T, val T) []T {
	for i, v := range slice {
		if v == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
