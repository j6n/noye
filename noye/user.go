package noye

import "fmt"

// User represents an IRC user
type User struct {
	Nick, User, Host string
}

func (u User) String() string {
	return fmt.Sprintf("%s!%s@%s", u.Nick, u.User, u.Host)
}
