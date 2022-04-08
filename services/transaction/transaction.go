package services

import (
	"context"
	"errors"
	"strconv"
	"time"
	"xendit/shared/constants"
	cs "xendit/shared/constants"
	"xendit/shared/utilities"
	utils "xendit/shared/utilities"
	vm "xendit/view_models/transaction"

	"github.com/go-redis/redis/v8"
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
	err := s.verifyAttemptThreshold(ctx, model)
	if err != nil {
		return &vm.TransactionResponse{
			Id:                "0",
			RequestData:       model,
			Recommendation:    cs.TransactionAccept,
			Context:           s.getTransactionContext(constants.TRANSACTION_STATE_FAILED, err),
			Created:           utils.UnixToISOTimeString(time.Now().UTC().UnixMilli()),
			Updated:           utils.UnixToISOTimeString(time.Now().UTC().UnixMilli()),
			TransactionStatus: s.getTransactionStatus(constants.TRANSACTION_STATE_FAILED),
		}, nil
	}
	// TODO: verify multiple unique cards
	err = s.verifyMultipleUniqueCards(ctx, model)
	if err != nil {
		return &vm.TransactionResponse{
			Id:                "0",
			RequestData:       model,
			Recommendation:    cs.TransactionAccept,
			Context:           s.getTransactionContext(constants.TRANSACTION_STATE_FAILED, err),
			Created:           utils.UnixToISOTimeString(time.Now().UTC().UnixMilli()),
			Updated:           utils.UnixToISOTimeString(time.Now().UTC().UnixMilli()),
			TransactionStatus: s.getTransactionStatus(constants.TRANSACTION_STATE_FAILED),
		}, nil
	}
	// TODO: verify conducted transaction
	err = s.verifyConductedTransaction(ctx, model)
	if err != nil {
		return &vm.TransactionResponse{
			Id:                "0",
			RequestData:       model,
			Recommendation:    cs.TransactionAccept,
			Context:           s.getTransactionContext(constants.TRANSACTION_STATE_FAILED, err),
			Created:           utils.UnixToISOTimeString(time.Now().UTC().UnixMilli()),
			Updated:           utils.UnixToISOTimeString(time.Now().UTC().UnixMilli()),
			TransactionStatus: s.getTransactionStatus(constants.TRANSACTION_STATE_FAILED),
		}, nil
	}

	newTransaction, err := s.repository.RecordTransaction(ctx, model)

	if err != nil || newTransaction == nil {
		return nil, err
	}

	err = s.addUniqueCardInUse(ctx, model)

	return &vm.TransactionResponse{
		Id:                strconv.FormatInt(newTransaction.Id, 10),
		RequestData:       model,
		Recommendation:    cs.TransactionAccept,
		Context:           s.getTransactionContext(newTransaction.State, err),
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

	if status {
		s.addConductedTransactionValue(ctx, data.Amount, data.CardID)
	} else {
		s.markTransactionFailedAttempt(ctx, data.CardID)
	}

	payload := struct {
		TransId string `json:"trans_id"`
	}{TransId: trans_id}

	return &vm.TransactionResponse{
		Id:                strconv.FormatInt(data.Id, 10),
		RequestData:       payload,
		Recommendation:    cs.TransactionAccept,
		Context:           s.getTransactionContext(data.State, nil),
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
		Context:           s.getTransactionContext(data.State, nil),
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

func (s *service) getTransactionContext(trans_state cs.TransactionState, err error) []string {
	if trans_state == cs.TRANSACTION_STATE_DEFAULT || trans_state == cs.TRANSACTION_STATE_SUCCESSED {
		return make([]string, 0)
	}
	//TODO: handle failed context
	if err != nil {
		return []string{err.Error()}
	}
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
	key := utilities.GetFullKey(constants.CKS_CARD_FAILED_ATTEMPT, model.CardID)
	val, err := s.cache.Get(ctx, key).Result()

	if err != nil || err == redis.Nil {
		s.logger.Infof("Key (%v) Empty result", key)
		return nil
	}

	count, _ := strconv.ParseInt(val, 10, 64)

	if count >= 3 {
		return errors.New(constants.AttemptThresholdBlocking)
	}
	return nil
}

// TODO: verify multiple unique cards
func (s *service) verifyMultipleUniqueCards(ctx context.Context, model *vm.TransactionRequest) error {
	key := utilities.GetFullKey(constants.CKS_USER_CARDS_IN_USE, model.UserID)
	vals, err := s.cache.SMembers(ctx, key).Result()

	if err != nil || err == redis.Nil {
		s.logger.Infof("Key (%v) Empty result", key)
		return nil
	}

	for _, v := range vals {
		if v == model.CardID && len(vals) >= 4 {
			return errors.New(constants.MultipleUniqueCardsBlocking)
		}
	}
	_, err = s.cache.SAdd(ctx, key, model.CardID).Result()
	s.cache.Do(ctx, "EXPIRE", key, 2*time.Minute)
	return err
}

// TODO: verify conducted transaction
func (s *service) verifyConductedTransaction(ctx context.Context, model *vm.TransactionRequest) error {
	key := utilities.GetFullKey(constants.CKS_USER_CONDUCTED_AMOUNT, model.UserID)
	val, err := s.cache.Get(ctx, key).Result()

	if err != nil || err == redis.Nil {
		s.logger.Infof("Key (%v) Empty result", key)
		return nil
	}

	amount, _ := strconv.ParseInt(val, 10, 64)

	if amount >= 1000000 {
		return errors.New(constants.ConductedExcessiveLoadingBlocking)
	}
	return nil
}

func (s *service) markTransactionFailedAttempt(ctx context.Context, cardId string) error {
	key := utilities.GetFullKey(cs.CKS_CARD_FAILED_ATTEMPT, cardId)
	oldVal, err := s.cache.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {
		_, err = s.cache.Set(ctx, key, 1, 1*time.Minute).Result()
		return err
	}

	old, err := strconv.ParseInt(oldVal, 10, 64)
	_, err = s.cache.Set(ctx, key, old+1, 1*time.Minute).Result()
	return nil
}

// add new card to current processing list
func (s *service) addUniqueCardInUse(ctx context.Context, model *vm.TransactionRequest) error {
	key := utilities.GetFullKey(cs.CKS_USER_CARDS_IN_USE, model.UserID)
	_, err := s.cache.SAdd(ctx, key, model.CardID).Result()
	if err != nil {
		return err
	}
	s.cache.Do(ctx, "EXPIRE", key, 2*time.Minute)
	return nil
}

// add conducted value to update total spending
func (s *service) addConductedTransactionValue(ctx context.Context, amount float64, cardId string) error {
	key := utilities.GetFullKey(cs.CKS_USER_CONDUCTED_AMOUNT, cardId)
	oldVal, err := s.cache.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {
		_, err = s.cache.Set(ctx, key, amount, 1*time.Minute).Result()
		return err
	}

	old, err := strconv.ParseFloat(oldVal, 64)
	_, err = s.cache.Set(ctx, key, old+amount, 1*time.Minute).Result()
	return nil
}
