package htmltest

import "net/http"

func statusCodeValid(code int) bool {
	return code == http.StatusPartialContent || code == http.StatusOK
}
