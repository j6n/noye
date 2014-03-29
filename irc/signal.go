package irc

import "sync"

type Signal struct {
	sig  chan struct{}
	once sync.Once
}

func NewSignal() *Signal {
	return &Signal{sig: make(chan struct{})}
}

func (s *Signal) Wait() <-chan struct{} {
	return s.sig
}

func (s *Signal) Close() {
	s.once.Do(func() { close(s.sig) })
}
