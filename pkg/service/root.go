package service

import "net/http"

func (s *Service) Root(w http.ResponseWriter, r *http.Request) {
	// data, err := ioutil.ReadFile("static/html/index.html")
	data, err := s.staticFS.ReadFile("static/html/index.html")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
