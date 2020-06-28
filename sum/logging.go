package sum

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"time"
)

func logingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger.Log("msg", "func log middleware start")
			defer logger.Log("msg", "func log middleware end")
			return e(ctx, request)
		}
	}

}

type loggingMiddleware struct {
	logger log.Logger
	next   SumService
}

func (mw loggingMiddleware) Sum(ctx context.Context, request *SumRequest) (resp *SumResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "sum",
			"input", fmt.Sprintf("%+v", request),
			"result", fmt.Sprintf("%+v", resp),
			"took", time.Since(begin),
		)
	}(time.Now())

	resp, err = mw.next.Sum(ctx, request)
	return
}
