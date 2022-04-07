package services

import (
	"context"
	config "xendit/config"
	itc "xendit/internal/interceptors"
	lg "xendit/pkg/logger"
	rps "xendit/repositories/transaction"
	vm "xendit/view_models/transaction"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
)

type TransactionService interface {
	// record a transaction and verify it
	RecordTransaction(ctx context.Context, model *vm.TransactionRequest) (*vm.TransactionResponse, error)

	// update state of a transaction was marked as in progress state
	FeedbackTransaction(ctx context.Context, trans_id string, status bool) (*vm.TransactionResponse, error)

	// get a transaction detail base on transaction_id
	GetTransaction(ctx context.Context, trans_id string) (*vm.TransactionResponse, error)

	// verify if a transaction is already exists
	IsTransactionExists(ctx context.Context, trans_id string) (bool, error)
}

type service struct {
	logger     lg.Logger
	tracer     opentracing.Tracer
	cfg        *config.Config
	cache      *redis.Client
	repository rps.TransactionRepository
	jwtManager *itc.JWTManager
}

// NewService func initializes a service
func NewService(logger lg.Logger,
	repository rps.TransactionRepository,
	trace opentracing.Tracer,
	jwtManager *itc.JWTManager,
	cfg *config.Config,
	cache *redis.Client,
) TransactionService {
	return &service{
		logger:     logger,
		cache:      cache,
		repository: repository,
		tracer:     trace,
		jwtManager: jwtManager,
		cfg:        cfg,
	}
}
