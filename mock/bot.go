package mock

// MockBot is a type that implements noye.Bot
// which can have parts of it overridden by changing the functors
type MockBot struct {
	SendFn    func(string, ...interface{})
	PrivmsgFn func(string, string)
	JoinFn    func(string)
	PartFn    func(string)
}

// NewMockBot returns a new MockBot with the fns set to no-ops
func NewMockBot() *MockBot {
	return &MockBot{
		func(string, ...interface{}) {},
		func(string, string) {},
		func(string) {},
		func(string) {},
	}
}

// Dial is just here for the interface
func (m *MockBot) Dial(addr, nick, user string) error { return nil }

// Close is just here for the interface
func (m *MockBot) Close() {}

// Send delegates to SendFn
func (m *MockBot) Send(f string, a ...interface{}) { m.SendFn(f, a) }

// Privmsg delegates to PrivmsgFn
func (m *MockBot) Privmsg(target, msg string) { m.PrivmsgFn(target, msg) }

// Join delegates to JoinFn
func (m *MockBot) Join(target string) { m.JoinFn(target) }

// Part delegates to PartFn
func (m *MockBot) Part(target string) { m.PartFn(target) }

// Wait is just here for the interface
func (m *MockBot) Wait() <-chan struct{} { return nil }

// Ready is just here for the interface
func (m *MockBot) Ready() <-chan struct{} { return nil }
