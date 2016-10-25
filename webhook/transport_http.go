package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func MakeWebhookHTTPServer(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	m := mux.NewRouter()
	m.Handle("/hook/{hookId}", httptransport.NewServer(
		ctx,
		endpoints.TriggerEndpoint,
		DecodeHTTPTriggerRequest,
		EncodeHTTPTriggerResponse,
		options...,
	))

	return m
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	if e, ok := err.(httptransport.Error); ok {
		msg = e.Err.Error()
		switch e.Domain {
		case httptransport.DomainDecode:
			code = http.StatusBadRequest

		case httptransport.DomainDo:
			code = http.StatusBadRequest
		}
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

func DecodeHTTPTriggerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	hookId := mux.Vars(r)["hookId"]
	return triggerRequest{hookId}, nil
}

func EncodeHTTPTriggerResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
