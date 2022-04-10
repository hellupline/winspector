package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hellupline/winspector/pkg/models"
	"github.com/hellupline/winspector/pkg/responses"
)

func (s *Service) RecordCreate(w http.ResponseWriter, r *http.Request) {
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
	recordKey := uuid.New()
	now := time.Now()
	requestData := models.NewRequestData(r)
	record := models.NewRecord(binKey, recordKey, now, requestData)
	s.DataStore.InsertRecord(record)
	w.WriteHeader(http.StatusCreated)
	go s.wsBroadcast(bin.BinKey, record)
}

func (s *Service) wsBroadcast(binKey uuid.UUID, record models.Record) {
	sockets, ok := s.DataStore.GetBinWatchers(binKey)
	if !ok {
		return
	}
	for _, conn := range sockets {
		writer, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println(err)
			s.DataStore.RemoveBinWatcher(binKey, conn)
			conn.Close()
			continue
		}
		response := responses.NewRecordResponse(record)
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			log.Println(err)
			s.DataStore.RemoveBinWatcher(binKey, conn)
			conn.Close()
			continue
		}
		if err := writer.Close(); err != nil {
			log.Println(err)
			s.DataStore.RemoveBinWatcher(binKey, conn)
			conn.Close()
			continue
		}
	}
}
