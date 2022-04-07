package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func binWatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	binKey, err := uuid.Parse(vars["binKey"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bin, ok := binStore[binKey]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	watchStore, ok := binWatchStore[bin.binKey]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	watchStore[conn] = true

	go connWatcher(conn, watchStore)
}

func connWatcher(conn *websocket.Conn, watchStore map[*websocket.Conn]bool) {
	defer func() {
		delete(watchStore, conn)
		conn.Close()
	}()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}
