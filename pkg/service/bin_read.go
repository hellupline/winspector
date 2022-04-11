package service

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hellupline/winspector/pkg/responses"
)

func (s *Service) BinRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	binKey, err := uuid.Parse(vars["binKey"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bin, ok := s.dataStore.GetBin(binKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	records, ok := s.dataStore.GetRecords(binKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response := responses.NewBinResponse(bin, records)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
