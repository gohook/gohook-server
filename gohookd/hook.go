package gohookd

type HookID string

type HookList []*Hook

type Hook struct {
	Id     HookID `json:"id"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

type HookCall struct {
	Id     HookID `json:"id"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

type HookRequest struct {
	Method string `json:"method"'`
}
