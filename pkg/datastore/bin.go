package datastore

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hellupline/winspector/pkg/models"
)

func (d *DataStore) InsertBin(bin models.Bin) {
	d.Lock()
	defer d.Unlock()
	d.BinStore[bin.BinKey] = bin
	d.BinRecordStore[bin.BinKey] = map[uuid.UUID]models.Record{}
	d.BinWatcherStore[bin.BinKey] = map[*websocket.Conn]bool{}
}

func (d *DataStore) GetBin(binKey uuid.UUID) (models.Bin, bool) {
	d.RLock()
	defer d.RUnlock()
	bin, ok := d.BinStore[binKey]
	return bin, ok
}
