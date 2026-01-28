package domain

type NtfyRequest struct {
	Title   string
	Message string
	Click   string
	Tag     string
	Actions []NtfyAction
}

type NtfyAction struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}
