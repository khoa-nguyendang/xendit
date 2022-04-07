package repositories

import (
	"context"
	"xendit/config"
	ent "xendit/entities"
	"xendit/pkg/logger"
	vm "xendit/view_models/login"

	"github.com/jmoiron/sqlx"
)

type authRepo struct {
	db     *sqlx.DB
	logger logger.Logger
	cfg    *config.Config
}

type AuthenticationRepository interface {
	GetUser(ctx context.Context, model *vm.LogInRequest) (*ent.User, error)
}

// NewRepository func initializes a service
func NewRepository(
	db *sqlx.DB,
	logger logger.Logger,

	cfg *config.Config,
) AuthenticationRepository {
	return &authRepo{db: db, logger: logger, cfg: cfg}
}
