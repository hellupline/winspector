package datastore

import (
	"sync"

	"github.com/google/uuid"
	"github.com/hellupline/winspector/pkg/models"
)

type DataStore struct {
	BinStore       map[uuid.UUID]models.Bin
	BinRecordStore map[uuid.UUID]map[uuid.UUID]models.Record

	// // using a single lock to avoit race condition between multiple locks
	sync.RWMutex
}

func NewDataStore() *DataStore {
	return &DataStore{
		BinStore:       map[uuid.UUID]models.Bin{},
		BinRecordStore: map[uuid.UUID]map[uuid.UUID]models.Record{},
	}
}
