package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func binRecordReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	binKey, err := uuid.Parse(vars["binKey"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	recordKey, err := uuid.Parse(vars["recordKey"])
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
	record, ok := recordStore[recordKey]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	headers := make(PairResponseList, 0, len(record.recordData.headers))
	for _, header := range record.recordData.headers {
		pairResponse := PairResponse{
			Key:   header.key,
			Value: header.value,
		}
		headers = append(headers, pairResponse)
	}
	response := RecordResponse{
		BinKey:    record.binKey.String(),
		RecordKey: record.recordKey.String(),
		CreatedAt: record.createdAt.Format(time.RFC3339),
		Method:    record.recordData.method,
		URL:       record.recordData.uRL,
		Headers:   headers,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
