package animations

import (
	"strconv"
	"strings"
)

type Animations struct {
	ID           int
	TextsID      int
	FileID       string
	MediaGroupID string
}

func (p *Animations) String() string {
	b := &strings.Builder{}

	b.WriteString("Animations{ID: ")
	b.WriteString(strconv.Itoa(p.ID))
	b.WriteString(", TextsID: ")
	b.WriteString(strconv.Itoa(p.TextsID))
	b.WriteString(", FileID: ")
	b.WriteString(p.FileID)
	b.WriteString(", MediaGroupID: ")
	b.WriteString(p.MediaGroupID)
	b.WriteRune('}')

	return b.String()
}
