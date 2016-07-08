package server

import (
	"net/http"
	"fmt"
)

func Run(htmlDir string) {
	http.Handle("/", http.FileServer(http.Dir(htmlDir)))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}
