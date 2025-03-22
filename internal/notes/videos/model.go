package videos

import (
	"strconv"
	"strings"
)

type Video struct {
	ID           int
	TextsID      int
	FileID       string
	MediaGroupID string
}

func (p *Video) String() string {
	b := &strings.Builder{}

	b.WriteString("Video{ID: ")
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
