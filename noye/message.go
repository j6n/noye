package noye

// User represents an IRC person
type User struct {
	Nick, User, Host string
}

// Message is a simple message
type Message struct {
	Target, Text string
	From         User
}

// IrcMessage is a representation of a raw IRC message
type IrcMessage struct {
	Source  User
	Command string
	Args    []string
	Text    string
	Raw     string
}
