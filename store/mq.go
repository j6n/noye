package store

import "sync"

type Queue struct {
	clients map[string]map[int64]struct{}
	mapping map[int64]chan string
	id      int64
	mu      sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		clients: make(map[string]map[int64]struct{}),
		mapping: make(map[int64]chan string),
	}
}

func (s *Queue) Update(key, val string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if ids, ok := s.clients[key]; ok {
		for id := range ids {
			if ch, ok := s.mapping[id]; ok && ch != nil {
				ch <- val
			}
		}
	}
}

func (s *Queue) Subscribe(key string) (int64, chan string) {
	id, ch := s.next(), make(chan string)

	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.clients[key]; !ok {
		s.clients[key] = map[int64]struct{}{id: struct{}{}}
	} else {
		m[id] = struct{}{}
	}

	s.mapping[id] = ch
	return id, ch
}

func (s *Queue) Unsubscribe(id int64) {
	s.mu.RLock()
	if ch, ok := s.mapping[id]; ok {
		close(ch)
	}

	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, m := range s.clients {
		delete(m, id)
	}

	delete(s.mapping, id)
}

func (s *Queue) next() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.id++
	return s.id
}
