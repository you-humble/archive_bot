package folder

import (
	"strconv"
	"strings"
)

type Folder struct {
	ID     int
	UserID int64
	Name   string
}

func (f *Folder) String() string {
	b := &strings.Builder{}

	b.WriteString("Folder{ID: ")
	b.WriteString(strconv.Itoa(f.ID))
	b.WriteString(", UserID: ")
	b.WriteString(strconv.FormatInt(f.UserID, 10))
	b.WriteString(", Name: ")
	b.WriteString(f.Name)
	b.WriteRune('}')

	return b.String()
}
