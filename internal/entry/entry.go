package entry

import "time"

type Type string

const (
	Group        Type = "group"
	File         Type = "file"
	FolderBrowse Type = "folder_browse"
	FolderView   Type = "folder_view"
	History      Type = "history"
	Favorite     Type = "favorite"
	Command      Type = "command"
	SearchResult Type = "search_result"
)

type Action string

const (
	ActionOpenFile     Action = "open_file"
	ActionBrowseFolder Action = "browse_folder"
	ActionViewFolder   Action = "view_folder"
	ActionNavigate     Action = "navigate"
	ActionRunCommand   Action = "run_command"
)

type Entry struct {
	ID          string            `json:"id"`
	Type        Type              `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Icon        string            `json:"icon,omitempty"`
	Action      Action            `json:"action"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

func (e Entry) FilterValue() string { return e.Name + " " + e.Description }

type Stat struct {
	Entry     Entry     `json:"entry"`
	Count     int       `json:"count"`
	FirstSeen time.Time `json:"first_seen"`
	LastUsed  time.Time `json:"last_used"`
	Pinned    bool      `json:"pinned,omitempty"`
}
