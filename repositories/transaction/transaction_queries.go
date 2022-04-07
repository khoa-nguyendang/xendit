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
	InsertTransactionQuery = `CALL dbo.sp_transactions_insert(?, ?, ?, ?, ?, ?)`

	/* Input:
	pUserId nvarchar(250),
	pCardId nvarchar(250),
	pTransactionId nvarchar(250),
	pAmount decimal,
	pTimestamp bigint(64),
	pState bit
	*/
	UpdateTransactionQuery = `CALL dbo.sp_transactions_update(?, ?, ?, ?, ?, ?)`

	/* Input:
	pTransactionId nvarchar(250),
	pTimestamp bigint(64),
	pState bit
	*/
	FeedbackTransactionQuery = `CALL dbo.sp_transactions_feedback(?, ?, ?)`

	/* pTransactionId nvarchar(250) */
	IsTransactionExistsQuery = `CALL dbo.sp_transactions_check_exists(?)`

	/* Input:
		pTransactionId nvarchar(250)
	Ouput:
		id              BIGINT(64)
	    transaction_id  VARCHAR(250),
	    created         BIGINT(64),
	    last_modified   BIGINT(64),
	    user_id			VARCHAR(250),
	    state			SMALLINT(8),
	*/
	GetTransactionByIdQuery = `CALL dbo.sp_transactions_get_by_transaction_id(?)`

	/* Input:
	pTransactionId nvarchar(250)
	pPayload blob,
	pResponse blob,
	pTimestamp bigint(64))
	*/
	InsertTransactionRequestHistoryQuery = `CALL dbo.sp_request_histories_insert(?, ?, ?, ?)`
)
