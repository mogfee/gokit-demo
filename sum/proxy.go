package sum

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type proxymw struct {
	ctx  context.Context
	next SumService
	sum  endpoint.Endpoint
}

func (mw proxymw) Sum(ctx context.Context, request *SumRequest) (*SumResponse, error) {
	response, err := mw.sum(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("返回值错误")
	}
	resp := response.(SumResponse)
	return &resp, nil
}
func makeSumProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/sum"
	}
	return httptransport.NewClient(
		"GET",
		u,
		func(ctx context.Context, request *http.Request, i interface{}) error {

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(i)
			if err != nil {
				return err
			}

			request.Body = ioutil.NopCloser(&buf)
			return nil
		}, func(ctx context.Context, response2 *http.Response) (response interface{}, err error) {
			resp := SumResponse{}
			if err := json.NewDecoder(response2.Body).Decode(&resp); err != nil {
				return nil, err
			}

			return resp, nil
		},
	).Endpoint()
}

type ServiceMiddleware func(SumService) SumService

func ProxyMiddleware(ctx context.Context, instance string, logger log.Logger) ServiceMiddleware {
	if instance == "" {
		logger.Log("proxy_to", "none")
		return func(service SumService) SumService {
			return service
		}
	}
	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond
	)
	var (
		instanceList = split(instance)
		endpointer   sd.FixedEndpointer
	)
	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeSumProxy(ctx, instance)
		//e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	retry := lb.Retry(maxAttempts, maxTime, lb.NewRoundRobin(endpointer))
	return func(service SumService) SumService {
		return proxymw{ctx: ctx, next: service, sum: retry}
	}
}
func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
