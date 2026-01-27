package ports

type NtfyClient interface {
	SendNotification(topic string, message string) error
}
