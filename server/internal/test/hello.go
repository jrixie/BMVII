package test

import (
	"log"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("hello world!")); err != nil {
		log.Println(err)
	}
}
