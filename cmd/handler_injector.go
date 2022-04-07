package main

import (
	"sync"
	"xendit/pkg/logger"
	"xendit/pkg/mysql"
	"xendit/pkg/trace"
	auth_rp "xendit/repositories/authentication"
	trans_rp "xendit/repositories/transaction"
	auth_sv "xendit/services/authentication"
	trans_sv "xendit/services/transaction"

	ict "xendit/internal/interceptors"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

func (s *server) addLogger() {
	appLogger := logger.NewApiLogger(s.cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s",
		s.cfg.Server.AppVersion,
		s.cfg.Logger.Level,
		s.cfg.Server.Mode,
	)
	appLogger.Infof("Success parsed config: %#v", s.cfg.Server.AppVersion)
	appLogger.Infof("mysql config: %#v", s.cfg.Mysql)
	s.logger = appLogger
}

func (s *server) addRedis() {
	//Redis
	redisMasterDb := redis.NewClient(&redis.Options{
		Addr:     s.cfg.Redis.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//Redis replica
	redisReplicas := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	s.redisWrite = redisMasterDb
	s.redisRead = redisReplicas
}

func (s *server) addMySQL() {
	mysqlDb, err := mysql.NewMysqlDB(s.cfg)

	if err != nil {
		s.logger.Infof("MySql init: %s", err)
	} else {
		s.logger.Infof("MySql connected: %#v", mysqlDb.Stats())
	}

	s.dbWrite = mysqlDb
}

func (s *server) addAuthenticationService() {
	authRepo := auth_rp.NewRepository(s.dbWrite, s.logger, s.cfg)
	s.authService = auth_sv.NewService(s.logger, authRepo, s.tracer, s.jwtManager, s.cfg, s.redisWrite)
}

func (s *server) addTransactionService() {
	transRepo := trans_rp.NewRepository(s.dbWrite, s.logger, s.cfg)
	s.transService = trans_sv.NewService(s.logger, transRepo, s.tracer, s.jwtManager, s.cfg, s.redisWrite)
}

func (s *server) addTracer() {
	tracer, err := trace.New("xenditTracer", s.cfg.Jaeger.Host)
	if err != nil {
		s.logger.Infof("trace new error: %v", err)
	}
	s.tracer = tracer
}

func (s *server) addJwt() {
	// jwtManager
	s.jwtManager = ict.NewJWTManager(s.logger, s.cfg)
}

// Clear resource of MySQL, Redis Connection
func (s *server) ClearResource() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)
	go func(dbWrite *sqlx.DB, wg *sync.WaitGroup) {
		_ = dbWrite.Close()
		defer wg.Done()
	}(s.dbWrite, &waitGroup)

	go func(rd *redis.Client, wg *sync.WaitGroup) {
		_ = rd.Close()
		defer wg.Done()
	}(s.redisWrite, &waitGroup)

	go func(rd *redis.Client, wg *sync.WaitGroup) {
		_ = rd.Close()
		defer wg.Done()
	}(s.redisRead, &waitGroup)
	waitGroup.Wait()
}
