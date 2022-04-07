package main

import (
	"time"

	"github.com/google/uuid"
)

type BinResponse struct {
	BinKey    string             `json:"bin_key"`
	CreatedAt string             `json:"created_at"`
	Records   RecordResponseList `json:"records"`
}

func NewBinResponse(bin Bin, recordStore map[uuid.UUID]Record) BinResponse {
	recordResponseList := NewRecordResponseList(recordStore)
	binResponse := BinResponse{
		BinKey:    bin.binKey.String(),
		CreatedAt: bin.createdAt.Format(time.RFC3339),
		Records:   recordResponseList,
	}
	return binResponse
}
