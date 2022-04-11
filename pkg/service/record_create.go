package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	bin, ok := s.dataStore.GetBin(binKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	recordKey := uuid.New()
	now := time.Now()
	requestData := models.NewRequestData(r)
	record := models.NewRecord(binKey, recordKey, now, requestData)
	s.dataStore.InsertRecord(record)
	w.WriteHeader(http.StatusCreated)
	response := responses.NewRecordResponse(record)
	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}
	s.publish(data, bin.BinKey)
}
