package stats

import (
	"context"

	"google.golang.org/grpc"
)

// gRPC 拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	resp, err = handler(ctx, req)

	return resp, err
}
