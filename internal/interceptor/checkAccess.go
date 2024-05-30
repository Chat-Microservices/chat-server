package interceptor

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/semho/chat-microservices/chat-server/internal/utils"
	accessV1 "github.com/semho/chat-microservices/chat-server/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(accessClient accessV1.AccessV1Client) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if accessClient == nil {
			return nil, status.Errorf(codes.Unauthenticated, "missing access client")
		}

		//спаны для трейсинга
		span, ctx := opentracing.StartSpanFromContext(ctx, "get access")
		defer span.Finish()
		span.SetTag("access", info.FullMethod)

		token, err := utils.GetToken(ctx)
		if err != nil {
			return nil, err
		}

		if err = utils.CheckAccess(ctx, info.FullMethod, accessClient, token); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
