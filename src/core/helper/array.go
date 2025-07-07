package helper

func IsContain[S ~[]E, E comparable](items S, item E) bool {
	for i := range items {
		if items[i] == item {
			return true
		}
	}
	return false
}
