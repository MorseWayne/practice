package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"time"

	v1 "github.com/MorseWayne/grpc-demo/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalculatorSerer struct {
	v1.UnimplementedCalculatorServiceServer
}

// Unary: 单词请求-响应
func (server *CalculatorSerer) Add(ctx context.Context, req *v1.AddRequest) (*v1.AddResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	select {
	case <-ctx.Done():
		{
			return nil, status.Error(codes.Canceled, "client canceled")
		}
	default:
	}
	return &v1.AddResponse{Result: req.A + req.B}, nil
}

// Client Streaming: 客户端流，客户端批量上传数据
func (server *CalculatorSerer) SumStream(stream v1.CalculatorService_SumStreamServer) error {
	var sum int64
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&v1.AddResponse{Result: sum})
		}
		if err != nil {
			return status.Errorf(codes.Internal, "recv error: %v", err)
		}
		sum = req.A + req.B
	}
}

// Server Streaming: 服务器流， 服务器批量推送数据
func (server *CalculatorSerer) RangeAdd(req *v1.RangeRequest, stream v1.CalculatorService_RangeAddServer) error {
	if req.Start > req.End {
		return status.Errorf(codes.InvalidArgument, "start[%d] > end[%d]", req.Start, req.End)
	}
	for i := req.Start; i <= req.End; i++ {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "client canceled")
		default:
		}
		if err := stream.Send(&v1.AddResponse{Result: i}); err != nil {
			return status.Errorf(codes.Internal, "send err : %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

// Bidirectional: 双向流，客户端服务端交替发送数据，适用于实时聊天
func (server *CalculatorSerer) ChatAdd(stream v1.CalculatorService_ChatAddServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "recv err : %v", err)
		}
		if err := stream.Send(&v1.AddResponse{Result: req.A + req.B}); err != nil {
			return status.Errorf(codes.Internal, "send err : %v", err)
		}
	}
}

func UnrayLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	log.Printf("[UNARY] method = %s, cost = %s, err = %v", info.FullMethod, time.Since(start), err)
	return
}

func StreamLoggingInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, ss)
	log.Printf("[STREAM] method = %s, cost = %s, err = %v", info.FullMethod, time.Since(start), err)
	return err
}

func NewGrpcServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(UnrayLoggingInterceptor),
		grpc.StreamInterceptor(StreamLoggingInterceptor),
	)
	v1.RegisterCalculatorServiceServer(s, &CalculatorSerer{})
	return s
}

func Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("grpc server listening at: %s", addr)
	return NewGrpcServer().Serve(lis)
}
