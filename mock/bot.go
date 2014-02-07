package mock

type MockBot struct {
	SendFn    func(string, ...interface{})
	PrivmsgFn func(string, string)
	JoinFn    func(string)
	PartFn    func(string)
}

func NewMockBot() *MockBot {
	return &MockBot{
		func(string, ...interface{}) {},
		func(string, string) {},
		func(string) {},
		func(string) {},
	}
}

func (m *MockBot) Dial(addr, nick, user string) error {
	return nil
}

func (m *MockBot) Close() {}

func (m *MockBot) Send(f string, a ...interface{}) {
	m.SendFn(f, a)
}

func (m *MockBot) Privmsg(target, msg string) {
	m.PrivmsgFn(target, msg)
}

func (m *MockBot) Join(target string) {
	m.JoinFn(target)
}

func (m *MockBot) Part(target string) {
	m.PartFn(target)
}

func (m *MockBot) Wait() <-chan struct{} {
	return nil
}
