package services

import (
	"context"
	config "xendit/config"
	"xendit/internal/interceptors"
	lg "xendit/pkg/logger"
	rps "xendit/repositories/authentication"
	vm "xendit/view_models/login"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
)

type AuthenticationService interface {
	LogIn(ctx context.Context, model *vm.LogInRequest) (*vm.LogInResponse, error)
}

type service struct {
	logger     lg.Logger
	tracer     opentracing.Tracer
	cfg        *config.Config
	cache      *redis.Client
	repository rps.AuthenticationRepository
	jwtManager *interceptors.JWTManager
}

// NewService func initializes a service
func NewService(logger lg.Logger,
	repository rps.AuthenticationRepository,
	trace opentracing.Tracer,
	jwtManager *interceptors.JWTManager,
	cfg *config.Config,
	cache *redis.Client,
) AuthenticationService {
	return &service{
		logger:     logger,
		cache:      cache,
		repository: repository,
		tracer:     trace,
		jwtManager: jwtManager,
		cfg:        cfg,
	}
}
