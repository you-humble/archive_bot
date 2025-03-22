package buttons

const (
	Prefix       string = "btn_"
	Delimiter    string = ":"
	WithPhoto    string = "ph"
	CreateFolder string = Prefix + "create_folder"
	DeleteFolder string = Prefix + "delete_folder"
	MoveNote     string = Prefix + "1_move"
	UpdateNote   string = Prefix + "2_update"
	DeleteNote   string = Prefix + "3_delete"
	Folders      string = "📁📁📁"
)

const DefaultFolderName = "Прочее"

var CatalogueOptions = map[string]string{
	MoveNote:   "📤",
	DeleteNote: "🗑️",
}

var MenuOptions = map[string]string{
	CreateFolder: "✅",
	DeleteFolder: "❌",
}
