package texts

import (
	"strconv"
	"strings"
	"time"
)

type TextNote struct {
	ID           int
	UserID       int64
	FolderID     int
	Type         string
	Description  string
	MediaGroupID string
	CreatedAt    time.Time
}

func (tn *TextNote) String() string {
	b := &strings.Builder{}

	b.WriteString("TextNote{ID: ")
	b.WriteString(strconv.Itoa(tn.ID))
	b.WriteString(", UserID: ")
	b.WriteString(strconv.FormatInt(tn.UserID, 10))
	b.WriteString(", FolderID: ")
	b.WriteString(strconv.Itoa(tn.FolderID))
	b.WriteString(", Type: ")
	b.WriteString(tn.Type)
	b.WriteString(", Description: ")
	b.WriteString(tn.Description)
	b.WriteString(", MediaGroupID: ")
	b.WriteString(tn.MediaGroupID)
	b.WriteString(", CreatedAt: ")
	b.WriteString(tn.CreatedAt.String())
	b.WriteRune('}')

	return b.String()
}
