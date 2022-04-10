package datastore

import (
	"github.com/google/uuid"
	"github.com/hellupline/winspector/pkg/models"
)

func (d *DataStore) GetRecords(binKey uuid.UUID) ([]models.Record, bool) {
	d.RLock()
	defer d.RUnlock()
	recordStore, ok := d.BinRecordStore[binKey]
	if !ok {
		return nil, false
	}
	records := make([]models.Record, 0, len(recordStore))
	for _, record := range recordStore {
		records = append(records, record)
	}
	return records, true
}

func (d *DataStore) InsertRecord(record models.Record) {
	d.Lock()
	defer d.Unlock()
	recordStore, ok := d.BinRecordStore[record.BinKey]
	if !ok {
		return
	}
	recordStore[record.RecordKey] = record
}

func (d *DataStore) GetRecord(binKey, recordKey uuid.UUID) (models.Record, bool) {
	d.RLock()
	defer d.RUnlock()
	recordStore, ok := d.BinRecordStore[binKey]
	if !ok {
		return models.Record{}, false
	}
	record, ok := recordStore[recordKey]
	return record, ok
}
