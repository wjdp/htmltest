package htmltest

import "log"

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
