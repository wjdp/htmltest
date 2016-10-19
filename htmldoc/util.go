package htmldoc

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func invalidPrePostRune(r rune) bool {
	// Return true if character at start or end of URL shoudl be trimmed
	return r == '\n' || r == ' '
}
