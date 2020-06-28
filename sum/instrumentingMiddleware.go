package sum

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           SumService
}

func (mw instrumentingMiddleware) Sum(ctx context.Context, request *SumRequest) (resp *SumResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "sum", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	resp, err = mw.next.Sum(ctx, request)
	return
}
