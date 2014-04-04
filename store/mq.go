package store

import "sync"

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
		return
	}

	q.mu.RLock()
	defer q.mu.RUnlock()

	if ids, ok := q.clients[key]; ok {
		for id := range ids {
			if ch, ok := q.mapping[id]; ok && ch != nil {
				ch <- val
			}
		}
	}
}

func (q *Queue) Subscribe(key string, private bool) (int64, chan string) {
	if _, ok := q.private[key]; ok && !private {
		// can't broadcast that
		return 0, nil
	}

	id, ch := q.next(), make(chan string)

	q.mu.Lock()
	defer q.mu.Unlock()
	if m, ok := q.clients[key]; !ok {
		q.clients[key] = map[int64]struct{}{id: struct{}{}}
	} else {
		m[id] = struct{}{}
	}

	q.mapping[id] = ch
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
