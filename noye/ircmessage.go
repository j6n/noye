package noye

// IrcMessage is a representation of a raw IRC message
type IrcMessage struct {
	Source  User
	Command string
	Args    []string
	Text    string
	Raw     string
}
