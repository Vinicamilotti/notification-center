package domain

type NotificationDTO struct {
	Topic               string
	Title               string
	Message             string
	AditionalAttributes map[string]any
	Actions             []NotificationAction
}

type NotificationAction struct {
	Label  string
	Action string
}
