package interceptor

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/guilhermealegre/go-clean-arch-core-lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	errorTag     = "error"
	errorListTag = "error-list"
)

func createErrorMD(respErr error) (metadata.MD, error) {
	switch respErr := respErr.(type) {
	case errors.ErrorDetails:
		bytes, err := json.Marshal(respErr)
		if err != nil {
			return nil, err
		}
		return metadata.Pairs(errorTag, string(bytes)), nil
	case errors.ErrorDetailsList:
		bytes, err := json.Marshal(respErr)
		if err != nil {
			return nil, err
		}
		return metadata.Pairs(errorListTag, string(bytes)), nil
	default:
		return nil, nil
	}
}

func receiveErrorMD(md metadata.MD) (err error) {
	var value []string
	if value = md.Get(errorTag); len(value) > 0 {
		newError := errors.ErrorDetails{}
		if err = json.Unmarshal([]byte(value[0]), &newError); err != nil {
			return err
		}
		if newError != (errors.ErrorDetails{}) {
			return newError
		}
	}

	if value = md.Get(errorListTag); len(value) > 0 {
		newError := errors.ErrorDetailsList{}
		if err = json.Unmarshal([]byte(value[0]), &newError); err != nil {
			return err
		}
		if !reflect.DeepEqual(newError, errors.ErrorDetailsList{}) {
			return newError
		}
	}

	return nil
}

func ErrorServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, respErr := handler(ctx, req)
		if respErr == nil {
			// no handling needed
			return resp, nil
		}

		md, err := createErrorMD(respErr)
		if err != nil {
			return resp, err
		}

		if md != nil {
			if err = grpc.SendHeader(ctx, md); err != nil {
				return resp, err
			}
		}

		return resp, respErr
	}
}

func ErrorServerStreamingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		respErr := handler(srv, ss)
		if respErr == nil {
			// no handling needed
			return nil
		}

		md, err := createErrorMD(respErr)
		if err != nil {
			return err
		}

		if md != nil {
			if err = ss.SetHeader(md); err != nil {
				return err
			}
		}

		return respErr
	}
}

func ErrorClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var md metadata.MD

		opts = append(opts, grpc.Header(&md))
		respErr := invoker(ctx, method, req, reply, cc, opts...)

		if respErr != nil && md != nil {
			if err := receiveErrorMD(md); err != nil {
				return err
			}
		}

		return respErr
	}
}

func ErrorClientStreamingInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var md metadata.MD

		opts = append(opts, grpc.Header(&md))
		cs, respErr := streamer(ctx, desc, cc, method, opts...)

		if respErr != nil && md != nil {
			if err := receiveErrorMD(md); err != nil {
				return cs, err
			}
		}

		return cs, respErr
	}
}
