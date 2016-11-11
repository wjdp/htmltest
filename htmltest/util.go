package htmltest

import "net/http"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func statusCodeValid(code int) bool {
	return code == http.StatusPartialContent || code == http.StatusOK
}
