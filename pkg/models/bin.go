package models

import (
	"time"

	"github.com/google/uuid"
)

type Bin struct {
	BinKey    uuid.UUID
	CreatedAt time.Time
}

func NewBin(binKey uuid.UUID, createdAt time.Time) Bin {
	return Bin{
		BinKey:    binKey,
		CreatedAt: createdAt,
	}
}
