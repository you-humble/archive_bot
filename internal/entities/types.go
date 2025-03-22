package entities

type Type int

const (
	Unknown Type = iota
	Message
	Photo
	Audio
	Document
	Video
	Animation
	Voice
)

func (t Type) String() string {
	switch t {
	case Unknown:
		return "unknown"
	case Message:
		return "message"
	case Photo:
		return "photo"
	case Audio:
		return "audio"
	case Document:
		return "doc"
	case Video:
		return "video"
	case Animation:
		return "animation"
	case Voice:
		return "voice"
	}
	return "unknown"
}

func ParseType(typeStr string) Type {
	switch typeStr {
	case "message":
		return Message
	case "photo":
		return Photo
	case "audio":
		return Audio
	case "doc":
		return Document
	case "video":
		return Video
	case "animation":
		return Animation
	case "voice":
		return Voice
	}
	return Unknown
}
