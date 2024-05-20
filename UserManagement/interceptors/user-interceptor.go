package interceptors

import (
	"context"
	"fmt"
	user "go-userm/proto/user/generatedFiles"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// func AuthInterceptor(authServiceClient user.AuthServiceClient) grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		if info.FullMethod == "/proto.user.UserService/Authenticate" {
// 			return handler(ctx, req)
// 		}

// 		md, ok := metadata.FromIncomingContext(ctx)
// 		if !ok {
// 			return nil, status.Errorf(codes.Unauthenticated, "Metadata not provided")
// 		}
// 		tokens, ok := md["authorization"]
// 		fmt.Println("Tokeni tokeni ljudi ", tokens)
// 		if !ok || len(tokens) < 1 {
// 			return nil, status.Errorf(codes.Unauthenticated, "Authorization token not found")
// 		}
// 		fmt.Println("Tokeni tokeni ljudi 2. put", tokens)

// 		// Create a new context with a timeout for the authentication request
// 		authCtx, cancel := context.WithTimeout(ctx, time.Second*5)
// 		defer cancel()

// 		// Make the authentication request to the AuthService
// 		_, err := authServiceClient.Authenticate(authCtx, &emptypb.Empty{})
// 		fmt.Println("Tokeni tokeni ljudi 3. put", tokens)
// 		if err != nil {
// 			return nil, status.Errorf(codes.Unauthenticated, "Authentication failed: %v", err)
// 		}
// 		fmt.Println("Tokeni tokeni ljudi 4. put", tokens)

// 		return handler(ctx, req)
// 	}
// }

func AuthInterceptor(authServiceClient user.AuthServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/UserService/UserSignUp" || info.FullMethod == "/UserService/Authenticate" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Metadata not provided")
		}
		tokens, ok := md["authorization"]
		if !ok || len(tokens) < 1 {
			return nil, status.Errorf(codes.Unauthenticated, "Authorization token not found")
		}

		authCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		authCtx = metadata.AppendToOutgoingContext(authCtx, "authorization", tokens[0])

		authResponse, err := authServiceClient.Authenticate(authCtx, &emptypb.Empty{})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Authentication failed: %v", err)
		}

		newMD := metadata.New(map[string]string{
			"id":       fmt.Sprintf("%d", authResponse.Id),
			"username": authResponse.Username,
			"password": authResponse.Password,
			"role":     fmt.Sprintf("%d", authResponse.Role),
		})

		md = metadata.Join(md, newMD)
		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
}
