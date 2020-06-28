package sum

import (
	"context"
	"errors"
)

type SumService interface {
	Sum(ctx context.Context, request *SumRequest) (*SumResponse, error)
}

type SumRequest struct {
	A int64
	B int64
}
type SumResponse struct {
	Sum int64
}

type sumService struct {
}

func (s sumService) Sum(ctx context.Context, request *SumRequest) (*SumResponse, error) {
	if request.A == 0 || request.B == 0 {
		return &SumResponse{}, errors.New("request error")
	}
	return &SumResponse{Sum: request.A + request.B}, nil
}
