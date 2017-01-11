package htmldoc

// Return true if character at start or end of URL should be trimmed.
func invalidPrePostRune(r rune) bool {
	return r == '\n' || r == ' '
}
