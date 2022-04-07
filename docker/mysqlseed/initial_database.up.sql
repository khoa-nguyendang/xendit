CREATE DATABASE IF NOT EXISTS xendit
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

use xendit;

CREATE TABLE IF NOT EXISTS users
(
    id               BIGINT(64)       AUTO_INCREMENT,
    username         VARCHAR(11)      NOT NULL,
    password         VARCHAR(2500)    NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transaction_histories
(
    id              BIGINT(64)      AUTO_INCREMENT,
    transaction_id  VARCHAR(250)    NOT NULL,
    created         BIGINT(64)      NOT NULL,
    last_modified   BIGINT(64)      NOT NULL,
    user_id			VARCHAR(250)	NOT NULL,
    state			SMALLINT(8)	NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS request_histories
(
    id              BIGINT(64)      AUTO_INCREMENT,
    transaction_id  VARCHAR(250)    NOT NULL,
    payload         BLOB    NOT NULL,
    response        BLOB    NOT NULL,
    timestamp       BIGINT(64)      NOT NULL,
);


CREATE INDEX  idx_users_user_name ON users (username); 
CREATE INDEX  idx_transaction_histories_user_id ON transaction_histories(user_id); 
CREATE INDEX  idx_transaction_histories_transaction_id ON transaction_histories(transaction_id); 
CREATE INDEX  idx_transaction_histories_user_id_transaction_id ON transaction_histories (user_id, transaction_id); 