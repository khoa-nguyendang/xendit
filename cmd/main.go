package main

import (
	"log"
	"os"
	"xendit/config"
	"xendit/internal/interceptors"
	"xendit/pkg/logger"
	authsv "xendit/services/authentication"
	transv "xendit/services/transaction"

	_ "xendit/cmd/docs"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server interface {
	Run() error
	GetLogger() logger.Logger
	ClearResource()
}

type server struct {
	mux          *gin.Engine
	logger       logger.Logger
	tracer       opentracing.Tracer
	cfg          *config.Config
	redisWrite   *redis.Client
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
		mux:          gin.Default(),
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
		mux: gin.Default(),
		cfg: cfg,
	}
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /
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

func (s *server) Run() error {
	// add necessary services
	s.addLogger()
	s.addTracer()
	s.addJwt()
	s.addRedis()
	s.addMySQL()
	s.addAuthenticationService()
	s.addTransactionService()

	// add apis
	s.mux.GET(PING, s.Ping)
	s.mux.POST(TRANSACTION_RECORD, s.TransactionRecordHandler)
	s.mux.POST(TRANSACTION_FEEDBACK, s.TransactionFeedbackHandler)
	s.mux.GET(TRANSACTION_GET, s.TransactionGetHandler)
	s.mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.logger.Infof("App is listening %v--", s.cfg.Server.Port)
	return s.mux.Run(s.cfg.Server.Port)
}

func (s *server) GetLogger() logger.Logger {
	return s.logger
}
