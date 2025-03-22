package router

import (
	"archive_bot/internal/const/buttons"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandAndTest(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		title string
		input string
		want  []string
	}{
		{
			"empty", "", []string{"", ""},
		},
		{
			"move last note alias", moveLastNoteAlias + "folder name", []string{moveLastNote, "folder name"},
		},
		{
			"move last note alias", buttons.Folders, []string{folders, ""},
		},
		{
			"some text without /", "some text without /", []string{"", "some text without /"},
		},
		{
			"command and text", "/command some text", []string{"/command", "some text"},
		},
		{
			"many slashes in a row", "//////////some text", []string{"", "//////////some text"},
		},
		{
			"many slashes in different places", "/command //:some/link", []string{"/command", "//:some/link"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			command, text := commandAndText(tc.input)
			got := []string{command, text}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestParseButtonCallback(t *testing.T) {
	// key + "_" + strconv.Itoa(noteID) + "_" + strconv.Itoa(folderID)
	t.Parallel()
	testCases := []struct {
		title string
		input string
		want  []int
	}{
		{
			"empty", "", []int{0, 0},
		},
		{
			"too long", buttons.CreateFolder + buttons.Delimiter + "1" + buttons.Delimiter + "2" + buttons.Delimiter + "3", []int{1, 2},
		},
		{
			"too short", buttons.CreateFolder + buttons.Delimiter + "1", []int{0, 0},
		},
		{
			"error in 1", buttons.CreateFolder + buttons.Delimiter + "error" + buttons.Delimiter + "2", []int{0, 0},
		},
		{
			"error in 2", buttons.CreateFolder + buttons.Delimiter + "1" + buttons.Delimiter + "error", []int{0, 0},
		},
		{
			"good case", buttons.CreateFolder + buttons.Delimiter + "1" + buttons.Delimiter + "2", []int{1, 2},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			noteID, FolderID, _ := ParseButtonCallback(tc.input)
			got := []int{noteID, FolderID}
			assert.Equal(t, tc.want, got)
		})
	}
}
