package channel

import (
	"fmt"

	"github.com/Vinicamilotti/notification-center/integration/ntfy/application"
	"github.com/Vinicamilotti/notification-center/lib/http"
	"github.com/Vinicamilotti/notification-center/shared/config"
	"github.com/Vinicamilotti/notification-center/shared/domain"
)

type NtfyChannel struct {
	Client     http.Client
	Config     config.NotificationConfig
	NtfyFacade *application.NtfyFacade
}

func NewNtfyChannel(config config.NotificationConfig) *NtfyChannel {
	return &NtfyChannel{
		Client:     http.NewClient(config.Channel),
		Config:     config,
		NtfyFacade: application.NewNtfyFacade(),
	}
}

func (nc *NtfyChannel) GetConfig() config.NotificationConfig {
	return nc.Config
}

func (nc *NtfyChannel) SendNotification(dto domain.NotificationDTO) error {
	ntfyReq := nc.NtfyFacade.ProcessRequest(dto)
	req, err := nc.Client.Post("", ntfyReq)

	if err != nil {
		return err
	}

	response, err := nc.Client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send ntfy notification, status code: %d", response.StatusCode)
	}

	return nil
}
