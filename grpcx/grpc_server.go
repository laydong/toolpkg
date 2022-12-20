package grpcx

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// GrpcServer struct
type GrpcServer struct {
	*grpc.Server
	opts   []grpc.UnaryServerInterceptor
	routes []func(server *GrpcServer)
}

// NewGrpcServer create new GrpcServer with default configuration
func NewGrpcServer() *GrpcServer {
	server := &GrpcServer{
		opts: []grpc.UnaryServerInterceptor{
			serverInterceptor,
		},
	}

	return server
}

func (gs *GrpcServer) Use(f ...func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)) {
	for _, vf := range f {
		gs.opts = append(gs.opts, vf)
	}
}

func (gs *GrpcServer) Register(f ...func(s *GrpcServer)) {
	gs.routes = append(gs.routes, f...)
}

func (gs *GrpcServer) Run(addr string) (err error) {
	// 初始化server, 将多个拦截器构建成一个拦截器
	gs.Server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(gs.opts...)),
	)

	// 注册路由
	for _, vf := range gs.routes {
		vf(gs)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	// 在给定的gRPC服务器上注册服务器反射服务
	reflection.Register(gs.Server)

	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们
	err = gs.Server.Serve(lis)
	return
}
