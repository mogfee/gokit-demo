package handlers

import (
	"context"

	pb "github.com/mogfee/gokit-demo/demo"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.SpiderServerServer {
	return spiderserverService{}
}

type spiderserverService struct{}

// ParseList implements Service.
func (s spiderserverService) ParseList(ctx context.Context, in *pb.ParseListRequest) (*pb.ParseListResponse, error) {
	var resp pb.ParseListResponse
	resp = pb.ParseListResponse{
		// Item:
	}
	return &resp, nil
}

// ParseDetail implements Service.
func (s spiderserverService) ParseDetail(ctx context.Context, in *pb.ParseDetailRequest) (*pb.ParseDetailResponse, error) {
	var resp pb.ParseDetailResponse
	resp = pb.ParseDetailResponse{
		// CompanyName:
		// Title:
		// City:
		// JobType:
		// Site:
		// BaseId:
		// Description:
		// Url:
		// JobCategory:
		// LastUpdateTime:
		// JobCreateTime:
		// Country:
		// Location:
	}
	return &resp, nil
}
