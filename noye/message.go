package noye

// Message is a simple message
type Message struct {
	From, Target, Text string
}

// IrcMessage is a representation of a raw IRC message
type IrcMessage struct {
	Source  string
	Command string
	Args    []string
	Text    string
	Raw     string
}
