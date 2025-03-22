package audios

import (
	"strconv"
	"strings"
)

type Audio struct {
	ID           int
	TextsID      int
	FileID       string
	MediaGroupID string
}

func (p *Audio) String() string {
	b := &strings.Builder{}

	b.WriteString("Audio{ID: ")
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
