package service

import (
	"embed"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hellupline/winspector/pkg/datastore"
	"golang.org/x/time/rate"
)

type Service struct {
	dataStore        *datastore.DataStore
	staticFS         embed.FS
	subscribers      map[uuid.UUID]map[*subscriber]struct{}
	subscribersMutex sync.Mutex
	publishLimiter   *rate.Limiter
}

func NewService(dataStore *datastore.DataStore, embedFS embed.FS) *Service {
	return &Service{
		dataStore:      dataStore,
		staticFS:       embedFS,
		subscribers:    map[uuid.UUID]map[*subscriber]struct{}{},
		publishLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router().ServeHTTP(w, r)
}

func (s *Service) Router() *mux.Router {
	// staticFileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static/")))
	staticFileServer := http.FileServer(http.FS(s.staticFS))
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
