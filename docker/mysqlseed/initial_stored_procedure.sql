use xendit;

delimiter $$
CREATE PROCEDURE `sp_transactions_insert` (
    pUserId nvarchar(250), 
    pCardId nvarchar(250), 
    pTransactionId nvarchar(250), 
    pAmount decimal,
    pTimestamp bigint(64),
    pState smallint(8))
BEGIN
	IF (SELECT EXISTS (SELECT id FROM transaction_histories WHERE transaction_id = pTransactionId) = 0) THEN
		INSERT INTO transaction_histories(user_id, card_id, transaction_id, amount, created, last_modified, state)
		VALUES (pUserId, pCardId, pTransactionId, pAmount, pTimestamp, pTimestamp, pState);
    END IF;
END $$
delimiter ;

delimiter $$
CREATE PROCEDURE `sp_transactions_update` (
    pUserId nvarchar(250), 
    pCardId nvarchar(250), 
    pTransactionId nvarchar(250), 
    pAmount decimal,
    pTimestamp bigint(64),
    pState smallint(8))
BEGIN
	UPDATE transaction_histories 
    SET user_id = pUserId, card_id = pCardId, amount = pAmount, last_modified = pTimestamp, state = pState
    WHERE  transaction_id = pTransactionId;
END $$
delimiter ;

delimiter $$
CREATE PROCEDURE `sp_transactions_feedback` (
    pTransactionId nvarchar(250), 
    pTimestamp bigint(64),
    pState smallint(8))
BEGIN
	UPDATE transaction_histories 
    SET last_modified = pTimestamp, state = pState
    WHERE  transaction_id = pTransactionId;
END $$
delimiter ;

delimiter $$
CREATE PROCEDURE `sp_request_histories_insert` (
    pTransactionId nvarchar(250), 
    pPayload blob,
    pResponse blob,
    pTimestamp bigint(64))
BEGIN
	UPDATE transaction_histories 
    SET last_modified = pTimestamp, state = pState
    WHERE  transaction_id = pTransactionId;
END $$
delimiter ;