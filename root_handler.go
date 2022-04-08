package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// data, err := ioutil.ReadFile("static/html/index.html")
	data, err := staticFS.ReadFile("static/html/index.html")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
