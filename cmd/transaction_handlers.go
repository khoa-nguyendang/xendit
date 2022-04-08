package main

import (
	"xendit/pkg/logger"

	_ "xendit/view_models/login"
	_ "xendit/view_models/transaction"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// GetUser godoc
// @Summary Record a transaction
// @Produce json
// @Body {object} viewmodels.TransactionRequest
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /transactions [post]
func (s *server) TransactionRecordHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "TransactionRecordHandler",
	})
}

// GetUser godoc
// @Summary Retrieves a transaction by transaction_id
// @Produce json
// @Param transaction_id path string true "transaction_id"
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /transactions/{transaction_id} [get]
func (s *server) TransactionGetHandler(c *gin.Context) {
	transaction_id := c.Param("transaction_id")
	c.JSON(200, gin.H{
		"message":        "TransactionGetHandler",
		"transaction_id": transaction_id,
	})
}

// GetUser godoc
// @Summary update feedback from payment gateway
// @Produce json
// @Param transaction_id path string true "transaction_id"
// @Body {object} viewmodels.TransactionRequest
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /transactions/{transaction_id}/feedback [post]
func (s *server) TransactionFeedbackHandler(c *gin.Context) {
	transaction_id := c.Param("transaction_id")
	c.JSON(200, gin.H{
		"message":        "TransactionFeedbackHandler",
		"transaction_id": transaction_id,
	})
}

// GetUser godoc
// @Summary update feedback from payment gateway
// @Produce json
// @Param transaction_id path string true "transaction_id"
// @Body {object} viewmodels.TransactionResponse
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /ping [get]
func (s *server) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
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

	s.logger.Debugf("App is listening %v--", s.cfg.Server.Port)
	return s.mux.Run(s.cfg.Server.Port)
}

func (s *server) GetLogger() logger.Logger {
	return s.logger
}
