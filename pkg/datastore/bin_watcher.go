package datastore

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (d *DataStore) GetBinWatchers(binKey uuid.UUID) ([]*websocket.Conn, bool) {
	d.RLock()
	defer d.RUnlock()
	watcherStore, ok := d.BinWatcherStore[binKey]
	if !ok {
		return nil, false
	}
	watchers := make([]*websocket.Conn, 0, len(watcherStore))
	for watcher := range watcherStore {
		watchers = append(watchers, watcher)
	}
	return watchers, true
}

func (d *DataStore) InsertBinWatcher(binKey uuid.UUID, conn *websocket.Conn) {
	d.RLock()
	defer d.RUnlock()
	binWatcherStore, ok := d.BinWatcherStore[binKey]
	if !ok {
		return
	}
	binWatcherStore[conn] = true
}

func (d *DataStore) RemoveBinWatcher(binKey uuid.UUID, conn *websocket.Conn) {
	d.RLock()
	defer d.RUnlock()
	binWatcherStore, ok := d.BinWatcherStore[binKey]
	if !ok {
		return
	}
	delete(binWatcherStore, conn)
}
