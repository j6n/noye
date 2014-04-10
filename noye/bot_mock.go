package noye

// MockBot is a type that implements noye.Bot
// which can have parts of it overridden by changing the functors
type MockBot struct {
	DialFn  func(string, string, string) error
	CloseFn func()

	SendFn    func(string, ...interface{})
	PrivmsgFn func(string, string)

	JoinFn func(string)
	PartFn func(string)
	QuitFn func()

	WaitFn  func() <-chan struct{}
	ReadyFn func() <-chan struct{}

	ManagerFn func() Manager
}

// NewMockBot returns a new MockBot with the fns set to no-ops
func NewMockBot() *MockBot {
	return &MockBot{
		DialFn:  func(string, string, string) error { return nil },
		CloseFn: func() {},

		SendFn:    func(string, ...interface{}) {},
		PrivmsgFn: func(string, string) {},

		JoinFn: func(string) {},
		PartFn: func(string) {},
		QuitFn: func() {},

		WaitFn:  func() <-chan struct{} { return nil },
		ReadyFn: func() <-chan struct{} { return nil },

		ManagerFn: func() Manager { return nil },
	}
}

// Dial delegates to DialFn
func (m *MockBot) Dial(addr, nick, user string) error {
	return m.DialFn(addr, nick, user)
}

// Close delegates to CloseFn
func (m *MockBot) Close() {
	m.CloseFn()
}

// Send delegates to SendFn
func (m *MockBot) Send(f string, a ...interface{}) {
	m.SendFn(f, a)
}

// Privmsg delegates to PrivmsgFn
func (m *MockBot) Privmsg(target, msg string) {
	m.PrivmsgFn(target, msg)
}

// Join delegates to JoinFn
func (m *MockBot) Join(target string) {
	m.JoinFn(target)
}

// Part delegates to PartFn
func (m *MockBot) Part(target string) {
	m.PartFn(target)
}

// Quit delegates to QuitFn
func (m *MockBot) Quit() {
	m.QuitFn()
}

// Wait delegates to WaitFn
func (m *MockBot) Wait() <-chan struct{} {
	return m.WaitFn()
}

// Ready delegates to ReadyFn
func (m *MockBot) Ready() <-chan struct{} {
	return m.ReadyFn()
}

// Manager delegates to ManagerFn
func (m *MockBot) Manager() Manager {
	return m.ManagerFn()
}
