package service

import (
	"embed"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hellupline/winspector/pkg/datastore"
)

type Service struct {
	DataStore *datastore.DataStore
	StaticFS  embed.FS
}

func NewService(dataStore *datastore.DataStore, embedFS embed.FS) *Service {
	return &Service{
		DataStore: dataStore,
		StaticFS:  embedFS,
	}
}

func (s *Service) Router() *mux.Router {
	// staticFileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static/")))
	staticFileServer := http.FileServer(http.FS(s.StaticFS))
	r := mux.NewRouter().StrictSlash(true)
	r.Path("/bin").Methods(http.MethodPost).Name("bin-create").HandlerFunc(s.BinCreate)
	r.Path("/bin/{binKey}").Methods(http.MethodGet).Name("bin-read").HandlerFunc(s.BinRead)
	r.Path("/bin/{binKey}/watch").Methods(http.MethodGet).Name("bin-watch").HandlerFunc(s.BinWatch)
	r.Path("/bin/{binKey}/records/{recordKey}").Methods(http.MethodGet).Name("record-read").HandlerFunc(s.RecordRead)
	r.PathPrefix("/record/{binKey}").Name("record-create").HandlerFunc(s.RecordCreate)
	r.PathPrefix("/static").Name("static").Handler(staticFileServer)
	r.PathPrefix("/").Methods(http.MethodGet).Name("root").HandlerFunc(s.Root)
	return r
}
