package main

import (
	"time"

	"github.com/google/uuid"
)

type Bin struct {
	binKey    uuid.UUID
	createdAt time.Time
}

func NewBin(binKey uuid.UUID, createdAt time.Time) Bin {
	return Bin{
		binKey:    binKey,
		createdAt: createdAt,
	}
}
