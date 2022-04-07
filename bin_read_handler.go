package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func binReadHandler(w http.ResponseWriter, r *http.Request) {
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
	binResponse := NewBinResponse(bin, recordStore)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(binResponse)
}
