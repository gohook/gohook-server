package webhook

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

type Middleware func(Service) Service

func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("layer", "endpoint", "error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}

func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw serviceLoggingMiddleware) Trigger(ctx context.Context) (v *WebhookStatus, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Trigger",
			"layer", "service",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Trigger(ctx)
}
