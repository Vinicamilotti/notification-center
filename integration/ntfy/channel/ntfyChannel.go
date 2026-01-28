package channel

import (
	"fmt"
	"io"
	"strings"

	http2 "net/http"

	"github.com/Vinicamilotti/notification-center/integration/ntfy/application"
	entities "github.com/Vinicamilotti/notification-center/integration/ntfy/domain"
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

func (nc *NtfyChannel) parseTags(tags []string) string {
	return strings.Join(tags, ",")
}

func (nc *NtfyChannel) determinateAction(action []entities.NtfyAction) string {
	s := ""
	for _, act := range action {
		partial := fmt.Sprintf("%s, %s, %s", act.Type, act.Label, act.Url)
		if act.Clear {
			partial = fmt.Sprintf("%s, clear=true", partial)
		}
		if act.Method != "" {
			partial = fmt.Sprintf("%s, method=%s", partial, act.Method)
		}
		if act.Body != "" {
			partial = fmt.Sprintf("%s, body=%s", partial, act.Body)
		}
		s += partial + ";"
	}
	return strings.TrimRight(s, ";")
}

func (nc *NtfyChannel) setHeaders(req *http2.Request, ntfyReq entities.NtfyRequest) {
	req.Header.Set("Title", ntfyReq.Title)
	req.Header.Set("tags", nc.parseTags(ntfyReq.Tag))
	req.Header.Set("click", ntfyReq.Click)
	req.Header.Set("actions", nc.determinateAction(ntfyReq.Actions))
}

func (nc *NtfyChannel) SendNotification(dto domain.NotificationDTO) error {
	ntfyReq := nc.NtfyFacade.ProcessRequest(dto)
	req, err := nc.Client.Post(fmt.Sprintf("/%s", dto.Topic), ntfyReq.Message)
	nc.setHeaders(req, *ntfyReq)

	if err != nil {
		return err
	}

	response, err := nc.Client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		read, _ := io.ReadAll(response.Body)
		fmt.Println(string(read))
		return fmt.Errorf("failed to send ntfy notification, status code: %d", response.StatusCode)
	}

	return nil
}
