package output

// Pluralise: a dumb pluraliser that only works for English.
func Pluralise(n int, singular string, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
