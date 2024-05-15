package utils

import (
	"context"
	accessV1 "github.com/semho/chat-microservices/chat-server/pkg/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func GetToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	accessToken, ok := getTokenFromMetadata(md)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "missing token")
	}

	//accessToken, _, ok := getTokensFromMetadata(md)
	//if !ok {
	//	return "", status.Errorf(codes.Unauthenticated, "missing token")
	//}
	//
	return accessToken, nil
}

func getTokenFromMetadata(md metadata.MD) (string, bool) {
	if authHeaders, ok := md["x-access-token"]; ok && len(authHeaders) > 0 {
		if strings.HasPrefix(authHeaders[0], "Bearer ") {
			return strings.TrimPrefix(authHeaders[0], "Bearer "), true
		}
	}
	return "", false
}

//через authorization постмана
//func getTokenFromMetadata(md metadata.MD) (string, bool) {
//	if authHeaders, ok := md["authorization"]; ok && len(authHeaders) > 0 {
//		if strings.HasPrefix(authHeaders[0], "Bearer ") {
//			fmt.Println(strings.TrimPrefix(authHeaders[0], "Bearer "))
//
//			return strings.TrimPrefix(authHeaders[0], "Bearer "), true
//		}
//	}
//	return "", false
//}

//метаданные постмана по refreshToken и accessToken
//func getTokensFromMetadata(md metadata.MD) (accessToken, refreshToken string, ok bool) {
//	if authHeaders, ok := md["x-access-token"]; ok && len(authHeaders) > 0 {
//		if strings.HasPrefix(authHeaders[0], "Bearer ") {
//			accessToken = strings.TrimPrefix(authHeaders[0], "Bearer ")
//		} else {
//			accessToken = authHeaders[0]
//		}
//	}
//
//	if refreshHeaders, ok := md["x-refresh-token"]; ok && len(refreshHeaders) > 0 {
//		if strings.HasPrefix(refreshHeaders[0], "Bearer ") {
//			refreshToken = strings.TrimPrefix(refreshHeaders[0], "Bearer ")
//		} else {
//			refreshToken = refreshHeaders[0]
//		}
//	}
//
//	// Если хотя бы один токен был найден, возвращаем true
//	ok = accessToken != "" || refreshToken != ""
//	return
//}

func CheckAccess(ctx context.Context, endpoint string, client accessV1.AccessV1Client, accessToken string) error {
	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := client.Check(ctx, &accessV1.CheckRequest{EndpointAddress: endpoint})
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "access denied: %v", err)
	}
	return nil
}
