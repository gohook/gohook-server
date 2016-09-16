package main

import (
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Endpoints struct {
	SumEndpoint endpoint.Endpoint
}

func (e Endpoints) Sum(ctx context.Context, a, b int) (int, error) {
	request := sumRequest{A: a, B: b}
	response, err := e.SumEndpoint(ctx, request)
	if err != nil {
		return 0, err
	}
	return response.(sumResponse).V, response.(sumResponse).Err
}

func MakeSumEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sumReq := request.(sumRequest)
		v, err := s.Sum(ctx, sumReq.A, sumReq.B)
		if err == ErrIntOverflow {
			return nil, err
		}
		return sumResponse{
			V:   v,
			Err: err,
		}, nil
	}
}

func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}

type sumRequest struct{ A, B int }

type sumResponse struct {
	V   int
	Err error
}
