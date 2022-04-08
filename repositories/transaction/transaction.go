package repositories

import (
	"context"
	"database/sql"
	"time"
	ent "xendit/entities"
	cs "xendit/shared/constants"
	vm "xendit/view_models/transaction"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func (r *transRepo) RecordTransaction(ctx context.Context, model *vm.TransactionRequest) (*ent.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TransactionRepository.RecordTransaction")
	defer span.Finish()

	r.logger.Infof("RecordTransaction:%v", model)
	if r.db == nil {
		r.logger.Info("RecordTransaction: Db is nil")
		return nil, errors.New("DB nil")
	}
	// real new transaction
	newTransaction, err := r.db.ExecContext(ctx, InsertTransactionQuery,
		model.UserID,
		model.CardID,
		model.TransactionID,
		model.Amount,
		time.Now().UTC().UnixMilli(),
		cs.TRANSACTION_STATE_DEFAULT,
	)

	if err != nil {
		r.logger.Errorf("RecordTransaction.error: %v", err)
		return nil, err
	}

	id, err := newTransaction.LastInsertId()

	if err != nil {
		r.logger.Errorf("RecordTransaction.error: %v", err)
		return nil, err
	}

	r.logger.Info("RecordTransaction.success: %v", id)
	return &ent.Transaction{
		Id:            id,
		UserID:        model.UserID,
		CardID:        model.CardID,
		TransactionID: model.TransactionID,
		Amount:        model.Amount,
	}, nil
}

// Update status of a transaction base on feedback from payment service
func (r *transRepo) FeedbackTransaction(ctx context.Context, trans_id string, status bool) (*ent.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TransactionRepository.LogIn")
	defer span.Finish()

	state := cs.TRANSACTION_STATE_FAILED
	if status {
		state = cs.TRANSACTION_STATE_SUCCESSED
	}
	_, err := r.db.ExecContext(ctx, FeedbackTransactionQuery,
		trans_id,
		time.Now().UTC().UnixMilli(),
		state,
	)

	if err != nil {
		r.logger.Errorf("FeedbackTransaction.error: %v", err)
		return nil, err
	}

	te, err := r.GetTransaction(ctx, trans_id)

	if err != nil {
		return nil, err
	}

	return te, nil
}

func (r *transRepo) GetTransaction(ctx context.Context, trans_id string) (*ent.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TransactionRepository.GetTransaction")
	defer span.Finish()

	result := &ent.Transaction{}
	if err := r.db.QueryRowContext(ctx, GetTransactionByIdQuery, trans_id).Scan(
		&result.Id,
		&result.TransactionID,
		&result.CardID,
		&result.Created,
		&result.LastModified,
		&result.UserID,
		&result.State,
		&result.State,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Errorf("GetTransaction.error: %v", err)
		return nil, err
	}

	return result, nil
}

func (r *transRepo) IsTransactionExists(ctx context.Context, trans_id string) (bool, error) {
	var transId int64 = 0
	if err := r.db.QueryRowContext(ctx, IsTransactionExistsQuery, trans_id).Scan(&trans_id); err != nil && errors.Cause(err) != sql.ErrNoRows {
		r.logger.Errorf("IsTransactionExists.error: %v", err)
		return true, err
	}

	return transId > 0, nil
}
