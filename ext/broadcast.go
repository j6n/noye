package ext

import "github.com/j6n/noye/store"

var mq = &store.Broadcast{
	Queue: store.NewQueue(),
	DB:    store.NewDB(),
}

// AddPrivate adds a list of internal keys to the message queue
func AddPrivate(keys ...string) { mq.Blacklist(keys...) }

// Broadcast broadcasts val to the internal key on the message queue
func Broadcast(key, val string) { mq.Update(key, val, true) }
