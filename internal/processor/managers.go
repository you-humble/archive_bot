package processor

import (
	"sync"

	"archive_bot/internal/entities"

	"github.com/looplab/fsm"
)

const (
	StartCreate  string = "start_create"
	SelectCreate string = "waiting_name_create"
	StartDelete  string = "start_delete"
	SelectDelete string = "waiting_name_delete"
	StartMove    string = "start_move"
	SelectMove   string = "waiting_name_move"
	StartUpdate  string = "start_update"
	SelectUpdate string = "waiting_name_update"
)

type folderManager struct {
	service FolderService

	mu               sync.RWMutex
	currentFolderIDs map[int64]int
	CreateStates     map[int64]*CreateState
	DeleteStates     map[int64]*DeleteState
}

func newFolderManager(folder FolderService) folderManager {
	return folderManager{
		service:          folder,
		mu:               sync.RWMutex{},
		currentFolderIDs: make(map[int64]int),
		CreateStates:     make(map[int64]*CreateState),
		DeleteStates:     make(map[int64]*DeleteState),
	}
}

func (fm *folderManager) CurrentFolderID(userID int64) int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.currentFolderIDs[userID]
}

func (fm *folderManager) SetCurrentFolderID(userID int64, folderID int) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.currentFolderIDs[userID] = folderID
}

type CreateState struct {
	FSM       *fsm.FSM
	MessageID int
}

type DeleteState struct {
	FSM       *fsm.FSM
	MessageID int
}

func (fm *folderManager) stateCreate(userID int64) *CreateState {
	fm.mu.RLock()
	state, ok := fm.CreateStates[userID]
	fm.mu.RUnlock()

	if !ok {
		return nil
	}
	return state
}

func (fm *folderManager) setStateCreate(userID int64) *CreateState {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	state := &CreateState{
		FSM: fsm.NewFSM(
			StartCreate,
			fsm.Events{
				{Name: "begin", Src: []string{StartCreate}, Dst: SelectCreate},
				{Name: "provide_name", Src: []string{SelectCreate}, Dst: StartCreate},
			},
			fsm.Callbacks{},
		),
	}
	fm.CreateStates[userID] = state
	return state
}

func (fm *folderManager) stateDelete(userID int64) *DeleteState {
	fm.mu.RLock()
	state, ok := fm.DeleteStates[userID]
	fm.mu.RUnlock()

	if !ok {
		return nil
	}
	return state
}

func (fm *folderManager) setStateDelete(userID int64) *DeleteState {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	state := &DeleteState{
		FSM: fsm.NewFSM(
			StartDelete,
			fsm.Events{
				{Name: "begin", Src: []string{StartDelete}, Dst: SelectDelete},
				{Name: "provide_name", Src: []string{SelectDelete}, Dst: StartDelete},
			},
			fsm.Callbacks{},
		),
	}
	fm.DeleteStates[userID] = state
	return state
}

type noteManager struct {
	texts     TextNoteService
	photos    PhotoNoteService
	documents DocsNoteService
	videos    VideoNoteService
	audios    AudioNoteService
	ani       AniNoteService
	voices    VoiceNoteService

	mu         sync.Mutex
	MoveStates map[int64]*MoveState
}

func newNoteManager(
	texts TextNoteService,
	photos PhotoNoteService,
	documents DocsNoteService,
	videos VideoNoteService,
	audios AudioNoteService,
	ani AniNoteService,
	voices VoiceNoteService,
) noteManager {
	return noteManager{
		texts:      texts,
		photos:     photos,
		documents:  documents,
		videos:     videos,
		audios:     audios,
		ani:        ani,
		voices:     voices,
		mu:         sync.Mutex{},
		MoveStates: make(map[int64]*MoveState),
	}
}

type MoveState struct {
	ParentFolderID int // TODO: return to the parent folder after the move
	NewFolderID    int
	NoteID         int
	FSM            *fsm.FSM
}

func (nm *noteManager) MoveState(userID int64, event *entities.Event) *MoveState {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	state, ok := nm.MoveStates[userID]

	if !ok {

		state = &MoveState{}
		state.FSM = fsm.NewFSM(
			StartMove,
			fsm.Events{
				{Name: "begin", Src: []string{StartMove}, Dst: SelectMove},
				{Name: "provide_ID", Src: []string{SelectMove}, Dst: StartMove},
			},
			fsm.Callbacks{},
		)
		nm.MoveStates[userID] = state
	}

	return state
}
