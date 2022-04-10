package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hellupline/winspector/pkg/models"
	"github.com/hellupline/winspector/pkg/responses"
)

func (s *Service) BinCreate(w http.ResponseWriter, r *http.Request) {
	binKey := uuid.New()
	now := time.Now()
	bin := models.NewBin(binKey, now)
	s.DataStore.InsertBin(bin)
	response := responses.NewBinResponse(bin, nil)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
