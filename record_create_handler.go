package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func recordCreateHandler(w http.ResponseWriter, r *http.Request) {
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
	recordStore, ok := binRecordStore[bin.binKey]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	watchStore, ok := binWatchStore[bin.binKey]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	recordKey := uuid.New()
	now := time.Now()
	recordData := NewRecordData(r)
	record := NewRecord(binKey, recordKey, now, recordData)
	recordStore[record.recordKey] = record
	w.WriteHeader(http.StatusCreated)

	{
		for conn := range watchStore {
			writer, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				conn.Close()
				delete(watchStore, conn)
				continue
			}
			recordResponse := NewRecordResponse(record)
			if err := json.NewEncoder(writer).Encode(recordResponse); err != nil {
				log.Println(err)
				conn.Close()
				delete(watchStore, conn)
				continue
			}
			if err := writer.Close(); err != nil {
				log.Println(err)
				conn.Close()
				delete(watchStore, conn)
				continue
			}
		}
	}
}
