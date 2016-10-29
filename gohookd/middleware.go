package gohookd

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gohook/gohook-server/user"
	"golang.org/x/net/context"
)

type Middleware func(Service) Service

func EndpointAuthMiddleware(logger log.Logger, auth user.AuthService) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			account, err := auth.AuthAccountFromToken(ctx.Value("token").(string))
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, "account", account)
			return next(ctx, request)
		}
	}
}

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

func (mw serviceLoggingMiddleware) List(ctx context.Context) (v HookList, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "List",
			"layer", "service",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.List(ctx)
}

func (mw serviceLoggingMiddleware) Create(ctx context.Context, request HookRequest) (v *Hook, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Create",
			"layer", "service",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Create(ctx, request)
}

func (mw serviceLoggingMiddleware) Delete(ctx context.Context, deleteID HookID) (v *Hook, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Delete",
			"layer", "service",
			"request", deleteID,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Delete(ctx, deleteID)
}
