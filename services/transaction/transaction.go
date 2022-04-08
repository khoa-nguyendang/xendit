package services

import (
	"context"
	"errors"
	"strconv"
	cs "xendit/shared/constants"
	utils "xendit/shared/utilities"
	vm "xendit/view_models/transaction"
)

func (s *service) RecordTransaction(ctx context.Context, model *vm.TransactionRequest) (*vm.TransactionResponse, error) {
	s.logger.Infof("RecordTransaction: %v", model)
	span := s.tracer.StartSpan("TransactionService.RecordTransaction")
	defer span.Finish()

	if model == nil {
		s.logger.Error("RecordTransaction Model Invalid")
		return nil, errors.New("RecordTransaction Model Invalid")
	}

	// TODO: verify attempt threshold

	// TODO: verify multiple unique cards

	// TODO: verify conducted transaction

	newTransaction, err := s.repository.RecordTransaction(ctx, model)

	if err != nil || newTransaction == nil {
		return nil, err
	}

	return &vm.TransactionResponse{
		Id:                strconv.FormatInt(newTransaction.Id, 10),
		RequestData:       model,
		Recommendation:    cs.TransactionAccept,
		Context:           s.getTransactionContext(newTransaction.State),
		Created:           utils.UnixToISOTimeString(newTransaction.Created),
		Updated:           utils.UnixToISOTimeString(newTransaction.LastModified),
		TransactionStatus: s.getTransactionStatus(newTransaction.State),
	}, nil
}

func (s *service) FeedbackTransaction(ctx context.Context, trans_id string, status bool) (*vm.TransactionResponse, error) {
	s.logger.Infof("RecordTransaction: %#v--%#v", trans_id, status)
	span := s.tracer.StartSpan("TransactionService.FeedbackTransaction")
	defer span.Finish()

	if trans_id == "" {
		s.logger.Info("FeedbackTransaction Model Invalid")
		return nil, errors.New("FeedbackTransaction Model Invalid")
	}

	_, err := s.repository.FeedbackTransaction(ctx, trans_id, status)

	if err != nil {
		return nil, err
	}

	data, err := s.repository.GetTransaction(ctx, trans_id)

	if err != nil || data == nil {
		return nil, err
	}

	payload := struct {
		TransId string `json:"trans_id"`
	}{TransId: trans_id}

	return &vm.TransactionResponse{
		Id:                strconv.FormatInt(data.Id, 10),
		RequestData:       payload,
		Recommendation:    cs.TransactionAccept,
		Context:           s.getTransactionContext(data.State),
		Created:           utils.UnixToISOTimeString(data.Created),
		Updated:           utils.UnixToISOTimeString(data.LastModified),
		TransactionStatus: s.getTransactionStatus(data.State),
	}, nil
}

func (s *service) GetTransaction(ctx context.Context, trans_id string) (*vm.TransactionResponse, error) {
	s.logger.Infof("RecordTransaction: %#v", trans_id)
	span := s.tracer.StartSpan("TransactionService.GetTransaction")
	defer span.Finish()

	if trans_id == "" {
		s.logger.Info("FeedbackTransaction Model Invalid")
		return nil, errors.New("FeedbackTransaction Model Invalid")
	}

	data, err := s.repository.GetTransaction(ctx, trans_id)

	if err != nil || data == nil {
		return nil, err
	}

	payload := struct {
		TransId string `json:"trans_id"`
	}{TransId: trans_id}
	return &vm.TransactionResponse{
		Id:                strconv.FormatInt(data.Id, 10),
		RequestData:       payload,
		Recommendation:    cs.TransactionAccept,
		Context:           s.getTransactionContext(data.State),
		Created:           utils.UnixToISOTimeString(data.Created),
		Updated:           utils.UnixToISOTimeString(data.LastModified),
		TransactionStatus: s.getTransactionStatus(data.State),
	}, nil

}

func (s *service) IsTransactionExists(ctx context.Context, trans_id string) (bool, error) {
	s.logger.Infof("RecordTransaction: %#v", trans_id)
	span := s.tracer.StartSpan("TransactionService.IsTransactionExists")
	defer span.Finish()

	if trans_id == "" {
		s.logger.Info("FeedbackTransaction Model Invalid")
		return false, errors.New("FeedbackTransaction Model Invalid")
	}

	data, err := s.repository.IsTransactionExists(ctx, trans_id)
	return data && err == nil, err
}

func (s *service) getTransactionContext(trans_state cs.TransactionState) []string {
	if trans_state == cs.TRANSACTION_STATE_DEFAULT || trans_state == cs.TRANSACTION_STATE_SUCCESSED {
		return make([]string, 0)
	}
	//TODO: handle failed context
	return make([]string, 0)
}

func (s *service) getTransactionStatus(trans_state cs.TransactionState) string {
	switch trans_state {
	case cs.TRANSACTION_STATE_DEFAULT:
		return "WAITING_FEEDBACK"
	case cs.TRANSACTION_STATE_FAILED:
		return "FAIL"
	case cs.TRANSACTION_STATE_SUCCESSED:
		return "SUCCESS"
	default:
		return "UNKNOWN"
	}
}

func (s *service) TransactionInsert(ctx context.Context, model *vm.TransactionRequest) (int64, error) {
	// There is no new record
	newTransaction, err := s.repository.RecordTransaction(ctx, model)

	if err != nil || newTransaction == nil {
		return 0, err
	}

	return newTransaction.Id, nil
}

// TODO: verify attempt threshold
func (s *service) verifyAttemptThreshold(ctx context.Context, model *vm.TransactionRequest) error {
	return nil
}

// TODO: verify multiple unique cards
func (s *service) verifyMultipleUniqueCards(ctx context.Context, model *vm.TransactionRequest) error {
	return nil
}

// TODO: verify conducted transaction
func (s *service) verifyConductedTransaction(ctx context.Context, model *vm.TransactionRequest) error {
	return nil
}

func (s *service) markTransactionFailedAttempt(ctx context.Context, model *vm.TransactionRequest) error {
	// key := utilities.GetFullKey(cs.CKS_CARD_FAILED_ATTEMPT, model.CardID)
	return nil
}
