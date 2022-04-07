package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func binCreateHandler(w http.ResponseWriter, r *http.Request) {
	binKey := uuid.New()
	now := time.Now()
	bin := NewBin(binKey, now)
	binStore[bin.binKey] = bin
	binRecordStore[bin.binKey] = map[uuid.UUID]Record{}
	binWatchStore[bin.binKey] = map[*websocket.Conn]bool{}
	response := BinResponse{
		BinKey:    bin.binKey.String(),
		CreatedAt: bin.createdAt.Format(time.RFC3339),
		Records:   RecordResponseList{},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
