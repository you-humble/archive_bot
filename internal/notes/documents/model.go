package documents

import (
	"strconv"
	"strings"
)

type Docs struct {
	ID           int
	TextsID      int
	FileID       string
	MediaGroupID string
}

func (d *Docs) String() string {
	b := &strings.Builder{}

	b.WriteString("Docs{ID: ")
	b.WriteString(strconv.Itoa(d.ID))
	b.WriteString(", TextsID: ")
	b.WriteString(strconv.Itoa(d.TextsID))
	b.WriteString(", FileID: ")
	b.WriteString(d.FileID)
	b.WriteString(", MediaGroupID: ")
	b.WriteString(d.MediaGroupID)
	b.WriteRune('}')

	return b.String()
}
