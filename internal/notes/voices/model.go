package voices

import (
	"strconv"
	"strings"
)

type Voice struct {
	ID      int
	TextsID int
	FileID  string
}

func (p *Voice) String() string {
	b := &strings.Builder{}

	b.WriteString("Voice{ID: ")
	b.WriteString(strconv.Itoa(p.ID))
	b.WriteString(", TextsID: ")
	b.WriteString(strconv.Itoa(p.TextsID))
	b.WriteString(", FileID: ")
	b.WriteString(p.FileID)
	b.WriteRune('}')

	return b.String()
}
