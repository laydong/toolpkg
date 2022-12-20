package grpcx

import (
	"context"
	"encoding/json"
	"github.com/laydong/toolpkg/metautils"
	"google.golang.org/grpc"
)

// serverInterceptor 提供服务的拦截器, 重写context, 记录出入参, 记录链路追踪
func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 初始化context
	md := metautils.ExtractIncoming(ctx)
	newCtx := NewGrpcContext(info.FullMethod, md)

	// 入参 header->meta
	reqByte, _ := json.Marshal(req)
	mdByte, _ := json.Marshal(md)
	newCtx.InfoF("%s", string(reqByte),
		newCtx.Field("header", string(mdByte)),
		newCtx.Field("path", info.FullMethod),
		newCtx.Field("protocol", protocol),
		newCtx.Field("title", "入参"))

	resp, err := handler(newCtx, req)

	respByte, _ := json.Marshal(resp)
	newCtx.InfoF("%s", string(respByte), newCtx.Field("title", "出参"))
	newCtx.SpanFinish(newCtx.TopSpan)
	return resp, err
}
