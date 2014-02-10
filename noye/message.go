package noye

type Message struct {
	From, Target, Text string
}

type IrcMessage struct {
	Source  string
	Command string
	Args    []string
	Text    string
}
