package store

var ConfigTable = "config"

// Broadcast wraps a queue and uses a db as a initial source
type Broadcast struct {
	*Queue
	DB *DB
}

// Init calls Subscribe, and then loads data from the shared table at key
// then it'll send the data over the channel
func (b *Broadcast) Init(table, key string, private bool) (int64, chan string) {
	id, ch := b.Queue.Subscribe(key, private)
	debug("'%d' subscribed to '%s' (%t)\n", id, key, private)
	data, err := b.DB.Get(ConfigTable, key)
	if err != nil {
		return id, ch
	}

	debug("'%d' sending data for '%s' (%t): %s\n", id, key, private, data)
	ch <- data

	data, err = Get(table, key)
	if err != nil {
		return id, ch
	}

	debug("'%d' sending data (script) for '%s' (%t): %s\n", id, key, private, data)
	ch <- data
	return id, ch
}
