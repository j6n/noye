package noye

// MockConn is a type that implements noye.Conn
// which can have parts of it overridden by changing the functors
type MockConn struct {
	DialFn      func(string) error
	CloseFn     func()
	WriteLineFn func(string)
	ReadLineFn  func() (string, error)
}

// NewMockConn returns a new MockConn with the fns set to no-ops
func NewMockConn() *MockConn {
	return &MockConn{
		DialFn:      func(string) error { return nil },
		CloseFn:     func() {},
		WriteLineFn: func(string) {},
		ReadLineFn:  func() (string, error) { return "", nil },
	}
}

// Dial delegates to DialFn
func (m *MockConn) Dial(addr string) error {
	return m.DialFn(addr)
}

// Close delegates to CloseFn
func (m *MockConn) Close() {
	m.CloseFn()
}

// WriteLine delegates to WriteLineFn
func (m *MockConn) WriteLine(raw string) {
	m.WriteLineFn(raw)
}

// ReadLine delegates to ReadLineFn
func (m *MockConn) ReadLine() (string, error) {
	return m.ReadLineFn()
}
