package repositories

const (
	/* Input:
	pUserId nvarchar(250),
	pCardId nvarchar(250),
	pTransactionId nvarchar(250),
	pAmount decimal,
	pTimestamp bigint(64),
	pState bit
	*/
	InsertTransactionQuery = `CALL sp_transactions_insert(?, ?, ?, ?, ?, ?)`

	/* Input:
	pUserId nvarchar(250),
	pCardId nvarchar(250),
	pTransactionId nvarchar(250),
	pAmount decimal,
	pTimestamp bigint(64),
	pState bit
	*/
	UpdateTransactionQuery = `CALL sp_transactions_update(?, ?, ?, ?, ?, ?)`

	/* Input:
	pTransactionId nvarchar(250),
	pTimestamp bigint(64),
	pState bit
	*/
	FeedbackTransactionQuery = `CALL sp_transactions_feedback(?, ?, ?)`

	/* pTransactionId nvarchar(250) */
	IsTransactionExistsQuery = `CALL sp_transactions_check_exists(?)`

	/* Input:
		pTransactionId nvarchar(250)
	Ouput:
		id              BIGINT(64)      AUTO_INCREMENT,
		transaction_id  VARCHAR(250)    NOT NULL,
		card_id         VARCHAR(250)    NOT NULL,
		created         BIGINT(64)      NOT NULL,
		last_modified   BIGINT(64)      NOT NULL,
		user_id			VARCHAR(250)	NOT NULL,
		amount			decimal	        NOT NULL,
		state			SMALLINT(8)	NOT NULL,
	*/
	GetTransactionByIdQuery = `CALL sp_transactions_get_by_transaction_id(?)`

	/* Input:
	pTransactionId nvarchar(250)
	pPayload blob,
	pResponse blob,
	pTimestamp bigint(64))
	*/
	InsertTransactionRequestHistoryQuery = `CALL sp_request_histories_insert(?, ?, ?, ?)`
)
