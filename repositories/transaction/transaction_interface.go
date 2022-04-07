package repositories

import (
	"context"
	"xendit/config"
	ent "xendit/entities"
	"xendit/pkg/logger"
	vm "xendit/view_models/transaction"

	"github.com/jmoiron/sqlx"
)

type transRepo struct {
	db     *sqlx.DB
	logger logger.Logger
	cfg    *config.Config
}

type TransactionRepository interface {
	RecordTransaction(ctx context.Context, model *vm.TransactionRequest) (*ent.Transaction, error)
	FeedbackTransaction(ctx context.Context, trans_id string, status bool) (*ent.Transaction, error)
	GetTransaction(ctx context.Context, trans_id string) (*ent.Transaction, error)
	IsTransactionExists(ctx context.Context, trans_id string) (bool, error)
}

// NewRepository func initializes a service
func NewRepository(
	db *sqlx.DB,
	logger logger.Logger,
	cfg *config.Config,
) TransactionRepository {
	return &transRepo{db: db, logger: logger, cfg: cfg}
}
