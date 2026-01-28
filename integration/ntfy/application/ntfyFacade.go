package application

import (
	"github.com/Vinicamilotti/notification-center/integration/ntfy/domain"
	entities "github.com/Vinicamilotti/notification-center/shared/domain"
)

type NtfyFacade struct{}

func NewNtfyFacade() *NtfyFacade {
	return &NtfyFacade{}
}

func (f *NtfyFacade) getAttribute(attributeName string, dto entities.NotificationDTO) string {
	val, exists := dto.AditionalAttributes[attributeName]
	if !exists {
		return ""
	}
	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal
}

func (f *NtfyFacade) determinateTag(statusValue string) string {
	switch statusValue {
	case "FIRING":
		return "rotating_light"
	case "RESOLVED":
		return "heavy_check_mark"
	default:
		return "warning"
	}
}

func (f *NtfyFacade) ProcessRequest(dto entities.NotificationDTO) *domain.NtfyRequest {

	ntfyReq := domain.NtfyRequest{
		Title:   dto.Title,
		Message: dto.Message,
		Click:   f.getAttribute("click", dto),
		Tag:     f.determinateTag(f.getAttribute("status", dto)),
		Actions: f.parseActions(dto),
	}

	return &ntfyReq
}

func getActionType(actionLabel entities.ActionType) domain.NtfyActionType {
	switch actionLabel {
	case entities.ActionTypeUrl:
		return domain.NtfyActionTypeView
	case entities.ActionTypeHttpCall:
		return domain.NtfyActionTypeHttp
	default:
		return domain.NtfyActionTypeView
	}
}

func (f *NtfyFacade) parseActions(dto entities.NotificationDTO) []domain.NtfyAction {
	actions := []domain.NtfyAction{}

	for _, action := range dto.Actions {
		ntfyAction := domain.NtfyAction{
			Type:  getActionType(action.Type),
			Label: action.Label,
			Url:   action.Action,
		}
		actions = append(actions, ntfyAction)
	}
	return actions
}
