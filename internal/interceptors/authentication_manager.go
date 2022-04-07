package interceptors

import (
	"context"
	"strconv"
	"xendit/config"
	ent "xendit/entities"
	"xendit/pkg/logger"
	"xendit/shared/constants"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthenticationInterceptor struct {
	logger     logger.Logger
	cfg        *config.Config
	jwtManager *JWTManager
}

func NewAuthenticationInterceptor(
	logger logger.Logger,
	cfg *config.Config,
	jwtManager *JWTManager) *AuthenticationInterceptor {
	return &AuthenticationInterceptor{
		logger,
		cfg,
		jwtManager,
	}
}

func (interceptor *AuthenticationInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		interceptor.logger.Infof("--> Unary interceptor: %s", info.FullMethod)
		new_ctx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			interceptor.logger.Errorf("Unary.Authentication %v", err)
			return nil, err
		}
		return handler(new_ctx, req)
	}
}

func (interceptor *AuthenticationInterceptor) authorize(ctx context.Context, method string) (newCtx context.Context, err error) {

	// Authentication logic
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.VerifyUser(accessToken)
	if err != nil || claims == nil {
		return ctx, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	return getAuthenticatedContext(ctx, claims), err
}

func getAuthenticatedContext(ctx context.Context, claims *ent.UserClaims) context.Context {
	customMetadatas := metadata.New(map[string]string{
		constants.UserId:   strconv.FormatInt(claims.Id, 10),
		constants.RoleId:   strconv.Itoa(int(claims.RoleId)),
		constants.Username: claims.Username,
	})
	return metadata.NewOutgoingContext(ctx, customMetadatas)
}
