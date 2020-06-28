package sum

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeSumEndpoint(srv SumService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		return srv.Sum(ctx, &req)
	}
}
