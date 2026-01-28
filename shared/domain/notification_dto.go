package domain

type ActionType string

const (
	ActionTypeUrl      ActionType = "url"
	ActionTypeHttpCall ActionType = "call"
)

type NotificationDTO struct {
	Topic               string
	Title               string
	Message             string
	AditionalAttributes map[string]any
	Actions             []NotificationAction
}

type NotificationAction struct {
	Type               ActionType
	Label              string
	Action             string
	AditionalParameter any
	Method             string
}
