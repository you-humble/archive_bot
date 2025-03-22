package user

import (
	"strconv"
	"strings"
)

type User struct {
	ID       int64
	Username string
}

func (u *User) String() string {
	b := &strings.Builder{}

	b.WriteString("User{ID: ")
	b.WriteString(strconv.FormatInt(u.ID, 10))
	b.WriteString(", Username: ")
	b.WriteString(u.Username)
	b.WriteRune('}')

	return b.String()
}
