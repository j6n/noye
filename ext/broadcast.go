package ext

import "github.com/j6n/noye/store"

var mq = &broadcaster{Queue: store.NewQueue()}

type broadcaster struct {
	*store.Queue
	m map[string]struct{}
}

// Init calls Subscribe, and then loads data from the shared table at key
// then it'll send the data over the channel
func (b *broadcaster) Init(table, key string, private bool) (int64, chan string) {
	id, ch := b.Queue.Subscribe(key, private)
	data, err := store.Get("shared", key)
	if err != nil {
		return id, ch
	}

	ch <- data

	data, err = store.Get(table, key)
	if err != nil {
		return id, ch
	}

	ch <- data
	return id, ch
}

// AddPrivate adds a list of internal keys to the message queue
func AddPrivate(keys ...string) { mq.Blacklist(keys...) }

// Broadcast broadcasts val to the internal key on the message queue
func Broadcast(key, val string) { mq.Update(key, val, true) }
