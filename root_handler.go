package main

import (
	"io/ioutil"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./index.html")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
