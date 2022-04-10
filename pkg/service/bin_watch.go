package service

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Service) BinWatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	binKey, err := uuid.Parse(vars["binKey"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bin, ok := s.DataStore.GetBin(binKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	s.DataStore.InsertBinWatcher(bin.BinKey, conn)
	go s.wsCleanUp(bin.BinKey, conn)
}

func (s *Service) wsCleanUp(binKey uuid.UUID, conn *websocket.Conn) {
	defer func() {
		s.DataStore.RemoveBinWatcher(binKey, conn)
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
