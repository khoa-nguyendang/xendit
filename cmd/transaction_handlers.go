package main

import (
	"context"
	"io/ioutil"
	"net/http"
	_ "xendit/view_models/login"
	vm "xendit/view_models/transaction"

	"github.com/gin-gonic/gin"
)

// GetUser godoc
// @Summary Record a transaction
// @Produce json
// @Param model body viewmodels.TransactionRequest true "Transactions"
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /transactions [post]
func (s *server) TransactionRecordHandler(c *gin.Context) {
	model := &vm.TransactionRequest{}
	c.BindJSON(model)
	if model == nil {
		s.sendBadRequestResponse(c)
		return
	}

	response, err := s.transService.RecordTransaction(context.Background(), model)
	if err != nil {
		s.logger.Infof("TransactionRecordHandler.err %v, %s", err.Error())
		s.sendCommonErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary Retrieves a transaction by transaction_id
// @Produce json
// @Param transaction_id path string true "transaction_id"
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /transactions/{transaction_id} [get]
func (s *server) TransactionGetHandler(c *gin.Context) {
	transaction_id := c.Param("transaction_id")
	if transaction_id == "" {
		s.sendBadRequestResponse(c)
		return
	}

	response, err := s.transService.GetTransaction(context.Background(), transaction_id)
	if err != nil {
		s.sendCommonErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary update feedback from payment gateway
// @Produce json
// @Param transaction_id path string true "transaction_id"
// @Param model body viewmodels.TransactionFeedbackRequest true "Feedback transaction"
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /transactions/{transaction_id}/feedback [post]
func (s *server) TransactionFeedbackHandler(c *gin.Context) {
	transaction_id := c.Param("transaction_id")
	model := &vm.TransactionFeedbackRequest{}
	c.Bind(&model)
	if model == nil {
		s.sendBadRequestResponse(c)
		return
	}

	response, err := s.transService.FeedbackTransaction(context.Background(), transaction_id, model.IsTransactionSuccess)
	if err != nil {
		s.sendCommonErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary update feedback from payment gateway
// @Produce json
// @Param transaction_id path string true "transaction_id"
// @Body {object} viewmodels.TransactionResponse
// @Success 200 {object} viewmodels.TransactionResponse
// @Router /ping [get]
func (s *server) Ping(c *gin.Context) {
	c.JSON(200, "pong")
}

func (s *server) sendBadRequestResponse(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"error":   "BadRequest",
		"payload": jsonData,
	})
}

func (s *server) sendForbiddenRequestResponse(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error
	}
	c.JSON(http.StatusForbidden, gin.H{
		"code":    http.StatusForbidden,
		"error":   "StatusForbidden",
		"payload": jsonData,
	})
}

func (s *server) sendCommonErrorResponse(c *gin.Context, err error) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error
	}
	message := "StatusExpectationFailed"
	if err != nil {
		message = err.Error()
	}
	c.JSON(http.StatusExpectationFailed, gin.H{
		"code":    http.StatusExpectationFailed,
		"error":   message,
		"payload": jsonData,
	})
}
