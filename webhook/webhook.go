package webhook

import (
	"github.com/gohook/gohook-server/gohookd"
)

type TriggerRequest struct {
	HookId gohookd.HookID
	Method string
	Body   []byte
}

type TriggerResponse struct {
	Code int `json:"code"`
}
