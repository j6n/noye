package script

import "github.com/j6n/noye/core/store"

var mq = &store.Broadcast{
	Queue: store.NewQueue(),
	DB:    store.NewDB(),
}

// AddPrivate adds a list of internal keys to the message queue
func AddPrivate(keys ...string) {
	mq.Blacklist(keys...)
}

// Broadcast the val to the internal key on the message queue
func Broadcast(key, val string) {
	mq.Update(key, val, true)
}
