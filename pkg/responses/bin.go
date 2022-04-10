package responses

import (
	"time"

	"github.com/hellupline/winspector/pkg/models"
)

type BinResponse struct {
	BinKey    string             `json:"bin_key"`
	CreatedAt string             `json:"created_at"`
	Records   RecordResponseList `json:"records"`
}

func NewBinResponse(bin models.Bin, records []models.Record) BinResponse {
	recordResponseList := NewRecordResponseList(records)
	return BinResponse{
		BinKey:    bin.BinKey.String(),
		CreatedAt: bin.CreatedAt.Format(time.RFC3339),
		Records:   recordResponseList,
	}
}
