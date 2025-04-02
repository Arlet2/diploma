package handlers

import "push_diploma/internal/core"

type PushPresenter struct {
	DeviceID string `json:"device_id"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

func (p PushPresenter) ToCore() core.Push {
	return core.Push{
		Title:    p.Title,
		Text:     p.Text,
		DeviceID: p.DeviceID,
	}
}

type SendResponsePresenter struct {
	PushID string `json:"push_id"`
}

type ErrorPresenter struct {
	Reason string `json:"reason"`
}
