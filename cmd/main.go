package main

import (
	"log"
	"net/http"
	"os"
	"xendit/config"
	"xendit/internal/interceptors"
	"xendit/pkg/logger"
	authsv "xendit/services/authentication"
	transv "xendit/services/transaction"

	"github.com/jmoiron/sqlx"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
)

type Server interface {
	Run() error
	GetLogger() logger.Logger
	ClearResource()
}

type server struct {
	mux          *http.ServeMux
	logger       logger.Logger
	tracer       opentracing.Tracer
	cfg          *config.Config
	redisWrite   *redis.Client
	redisRead    *redis.Client
	dbRead       *sqlx.DB
	dbWrite      *sqlx.DB
	authService  authsv.AuthenticationService
	transService transv.TransactionService
	jwtManager   *interceptors.JWTManager
}

func NewServer(
	logger logger.Logger,
	tracer opentracing.Tracer,
	cfg *config.Config,
	auth_sv authsv.AuthenticationService,
	trans_sv transv.TransactionService,
	jwtManager *interceptors.JWTManager,
) Server {
	return &server{
		mux:          http.NewServeMux(),
		logger:       logger,
		tracer:       tracer,
		jwtManager:   jwtManager,
		cfg:          cfg,
		authService:  auth_sv,
		transService: trans_sv,
	}
}

func NewBasicServer(cfg *config.Config) Server {
	return &server{
		mux: http.NewServeMux(),
		cfg: cfg,
	}
}

func main() {
	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)

	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}
	srv := NewBasicServer(cfg)

	defer srv.ClearResource()

	if err := srv.Run(); err != nil {
		srv.GetLogger().Infof("run got error: %v", err)
	}
}
