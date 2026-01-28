package domain

type NtfyActionType string

const (
	NtfyActionTypeView NtfyActionType = "view"
	NtfyActionTypeHttp NtfyActionType = "http"
)

type NtfyRequest struct {
	Title   string       `json:"title"`
	Message string       `json:"message"`
	Click   string       `json:"click"`
	Tag     []string     `json:"tag"`
	Actions []NtfyAction `json:"actions"`
}

type NtfyAction struct {
	Type  NtfyActionType `json:"action,omitempty"`
	Label string         `json:"label,omitempty"`
	Url   string         `json:"url,omitempty"`
	Clear bool           `json:"clear,omitempty"`
	Body  string         `json:"body,omitempty"`
}
