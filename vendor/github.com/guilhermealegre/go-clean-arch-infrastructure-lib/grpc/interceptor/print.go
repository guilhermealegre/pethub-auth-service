package interceptor

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
)

const (
	grpcServerMessage  = "server"
	grpcClientMessage  = "client"
	grpcSuccessMessage = "200"
	grpcErrorMessage   = "400"
	separator          = "|"
)

func status(success bool) string {
	message := grpcSuccessMessage
	if !success {
		message = grpcErrorMessage
	}
	return message
}

func printServer(logger io.Writer, method string, success bool) (n int, err error) {
	return logger.Write([]byte(fmt.Sprintln(grpcServerMessage, separator, status(success), separator, method)))
}

func printClient(logger io.Writer, method string, success bool) (n int, err error) {
	return logger.Write([]byte(fmt.Sprintln(grpcClientMessage, separator, status(success), separator, method)))
}

func PrintServerInterceptor(logger io.Writer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, respErr := handler(ctx, req)
		_, _ = printServer(logger, info.FullMethod, respErr == nil)
		return resp, respErr
	}
}

func PrintServerStreamingInterceptor(logger io.Writer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		respErr := handler(srv, ss)
		_, _ = printServer(logger, info.FullMethod, respErr == nil)
		return respErr
	}
}

func PrintClientInterceptor(logger io.Writer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		respErr := invoker(ctx, method, req, reply, cc, opts...)
		_, _ = printClient(logger, method, respErr == nil)
		return respErr
	}
}

func PrintClientStreamingInterceptor(logger io.Writer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		cs, respErr := streamer(ctx, desc, cc, method, opts...)
		_, _ = printClient(logger, method, respErr == nil)
		return cs, respErr
	}
}
