package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var binStore = map[uuid.UUID]Bin{}
var binRecordStore = map[uuid.UUID]map[uuid.UUID]Record{}
var binWatchStore = map[uuid.UUID]map[*websocket.Conn]bool{}

func main() {
	{
		binKey := uuid.MustParse("d45a2464-4bce-4628-95be-8b8dfebe90be")
		now := time.Now()
		bin := NewBin(binKey, now)
		binStore[bin.binKey] = bin
		binRecordStore[bin.binKey] = map[uuid.UUID]Record{}
		binWatchStore[bin.binKey] = map[*websocket.Conn]bool{}
	}

	staticFileServer := http.FileServer(http.Dir("./static/"))
	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", staticFileServer))
	r.Path("/").Methods("GET").HandlerFunc(rootHandler)
	r.Path("/bin").Methods("POST").HandlerFunc(binCreateHandler)
	r.Path("/bin/{binKey}").Methods("GET").HandlerFunc(binReadHandler)
	r.Path("/bin/{binKey}/watch").Methods("GET").HandlerFunc(binWatchHandler)
	r.Path("/bin/{binKey}/records/{recordKey}").Methods("GET").HandlerFunc(binRecordReadHandler)
	r.PathPrefix("/record/{binKey}").HandlerFunc(recordCreateHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}
