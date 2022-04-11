package service

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hellupline/winspector/pkg/responses"
)

func (s *Service) RecordRead(w http.ResponseWriter, r *http.Request) {
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
	bin, ok := s.dataStore.GetBin(binKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	record, ok := s.dataStore.GetRecord(bin.BinKey, recordKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response := responses.NewRecordResponse(record)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
