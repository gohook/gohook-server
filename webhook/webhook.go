package webhook

import (
	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/user"
)

type TriggerRequest struct {
	AccountId user.AccountId
	HookId    gohookd.HookID
	Method    string
	Body      []byte
}

type TriggerResponse struct {
	Code int `json:"code"`
}
