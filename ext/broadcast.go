package ext

// AddPrivate adds a list of internal keys to the message queue
func AddPrivate(keys ...string) { mq.Blacklist(keys...) }

// Broadcast broadcasts val to the internal key on the message queue
func Broadcast(key, val string) { mq.Update(key, val, true) }
