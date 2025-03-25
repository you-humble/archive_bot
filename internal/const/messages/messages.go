package messages

const (
	StartCommand string = `Привет! Я – твой личный помощник, который помогает создавать уютные папочки для заметок и ссылок.😬 
Со мной ты можешь:

• Организовывать свои идеи и интересные материалы в папки  📁
• Быстро находить сохранённые заметки и ссылки именно тогда, когда они понадобятся 🔍

Чтобы начать – просто отправь мне сообщение с заметкой или ссылкой – и я сразу её сохраню для тебя. Если хочешь увидеть свои папки или создать новую - нажми кнопку 📂📂📂.

Посмотреть инструкцию можно, набрав команду /info`
	FoldersCaption string = "📌 Мои папки"
	UnknownCommand string = "Не знаю о чем ты. О чем ты? 🤔"
	Error          string = "Ой! Не делай так... 😵"
	NotesIsEmpty   string = "Ничего нет 🕵🏼"
	NoteCreated    string = "Запись добавлена ✏️"
	Moved          string = "✏️"
	NoteRemoved    string = "Запись удалена 🧹"
	EmptyMessage   string = "▲▲▲▲▲"
	InfoVideo1     string = "BAACAgIAAxkBAAJKW2fgJe-M4SzvW1uvJNIP5JvaNO70AAJgawACzAwAAUtn4sBEP276AAE2BA"
	InfoVideo2     string = "BAACAgIAAxkBAAJKXGfgJe-ZHLN7SkOLTofpdDvEJsmEAAJhawACzAwAAUvPfV2u_yGyNzYE"
	InfoVideo3     string = "BAACAgIAAxkBAAJKXWfgJe-SHdLJCa8-bVcKSvPdP_eKAAJiawACzAwAAUsZ5QqXvx__eDYE"
	InfoVideo4     string = "BAACAgIAAxkBAAJKXmfgJe-o628i1tfcrmRCJSTVBJAFAAJkawACzAwAAUu7iI-JnZ0EqDYE"
	InfoVideo5     string = "BAACAgIAAxkBAAJKX2fgJe-2BRgqwzm_78tOq87a__ERAAJjawACzAwAAUsnF-m5Z4dWLjYE"
	InfoMessage1   string = "1. Чтобы добавить новую заметку, нужно написать или прислать что-то в бота.\nЧтобы стереть все кроме главного меню нажмите на 📂📂📂 или введите /folders"
	InfoMessage2   string = "2. Добавить новую папку можно, нажав на левую кнопку главного меню. Чтобы удалить папку, нужно нажать на правую кнопку."
	InfoMessage3   string = "3. Если написать или добавить что-то в бота, находясь в папке, новая запись сохранится в эту папку"
	InfoMessage4   string = "4. Левая кнопка под записью перемещает ее в нужную папку (после нажатия этой кнопки нужно выбрать папку, в которую необходимо переместить запись).\nПравая кнопка удаляет запись"
	InfoMessage5   string = "5. При добавлении записей через опцию 'Поделиться',\nесли в сообщении написать восклицательный знак и название папки (!название), то запись будеть добавлена в эту папку. Если этой папки не существует, она создастся автоматически."
)

var InfoMap = map[string]string{
	InfoVideo1: InfoMessage1,
	InfoVideo2: InfoMessage2,
	InfoVideo3: InfoMessage3,
	InfoVideo4: InfoMessage4,
	InfoVideo5: InfoMessage5,
}

const (
	FolderEmoji          string = "📁"
	AskFolderName        string = "Название папки?"
	FolderIsEmpty        string = "Нету папочек 🕵🏼"
	FolderCreated        string = "А вот и папочка ✏️"
	FolderNotExists      string = "Нет такой папки"
	FolderDeleted        string = "Папка удалена"
	WrongFolder          string = "Эту папку нельзя удалить"
	ChooseFolderToMove   string = "Выбери папку, в которую хочешь переместить заметку"
	ChooseFolderToDelete string = "Выбери папку, которую хочешь удалить"
)
