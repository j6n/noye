package store

import (
	"sync"

	"github.com/j6n/noye/logger"
)

var log = logger.Get()
var Debug bool

var debug = func(f string, a ...interface{}) {
	if Debug {
		log.Debugf(f, a...)
	}
}

type Queue struct {
	clients map[string]map[int64]struct{}
	mapping map[int64]chan string
	private map[string]struct{}

	id int64
	mu sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		clients: make(map[string]map[int64]struct{}),
		mapping: make(map[int64]chan string),
		private: make(map[string]struct{}),
	}
}

func (q *Queue) Blacklist(keys ...string) {
	val := struct{}{}
	for _, key := range keys {
		q.private[key] = val
	}
}

func (q *Queue) Update(key, val string, private bool) {
	if _, ok := q.private[key]; ok && !private {
		// can't broadcast that
		debug("can't broadcast on '%s' '%t'\n", key, private)
		return
	}

	q.mu.RLock()
	defer q.mu.RUnlock()

	var temp []int64
	if ids, ok := q.clients[key]; ok {
		for id := range ids {
			if ch, ok := q.mapping[id]; ok && ch != nil {
				temp = append(temp, id)
				go func(ch chan string, val string) {
					ch <- val
				}(ch, val)
			}
		}
	}

	debug("sending %s: '%s' to %v\n", key, val, temp)
}

func (q *Queue) Subscribe(key string, private bool) (int64, chan string) {
	if _, ok := q.private[key]; ok && !private {
		// can't listen to that
		debug("can't listen on '%s' '%t'\n", key, private)
		return 0, nil
	}

	id, ch := q.next(), make(chan string, 32)

	q.mu.Lock()
	defer q.mu.Unlock()
	if m, ok := q.clients[key]; !ok {
		q.clients[key] = map[int64]struct{}{id: struct{}{}}
	} else {
		m[id] = struct{}{}
	}

	q.mapping[id] = ch

	debug("'%d' subscribing to: %s (%t)\n", id, key, private)
	return id, ch
}

func (q *Queue) Unsubscribe(id int64) {
	q.mu.RLock()
	if ch, ok := q.mapping[id]; ok {
		close(ch)
	}

	q.mu.RUnlock()

	q.mu.Lock()
	defer q.mu.Unlock()
	for _, m := range q.clients {
		delete(m, id)
	}

	delete(q.mapping, id)
}

func (q *Queue) next() int64 {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.id++
	return q.id
}
