package irc

import "sync"

// Signal is a single-use blocking channel
type Signal struct {
	sig    chan struct{}
	once   sync.Once
	closed bool
}

// NewSignal returns a new Signal
func NewSignal() *Signal {
	return &Signal{sig: make(chan struct{})}
}

// Wait returns the internal channel
func (s *Signal) Wait() <-chan struct{} {
	return s.sig
}

// Close closes the channel, only once
func (s *Signal) Close() {
	s.once.Do(func() { close(s.sig); s.closed = true })
}

// Done returns whether the signal has been closed
func (s *Signal) Done() bool {
	return s.closed
}
