package models

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	BinKey      uuid.UUID
	RecordKey   uuid.UUID
	CreatedAt   time.Time
	RequestData RequestData
}

func NewRecord(binKey, recordKey uuid.UUID, createdAt time.Time, requestData RequestData) Record {
	return Record{
		BinKey:      binKey,
		RecordKey:   recordKey,
		CreatedAt:   createdAt,
		RequestData: requestData,
	}
}
