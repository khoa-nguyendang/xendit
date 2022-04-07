package constants

const (
	UserId   string = "user_id"
	Username string = "username"
	RoleId   string = "role_id"
)

const (
	TransactionAccept string = "ACCEPT"
	TransactionReject string = "REJECT"
)

const (
	AttemptThresholdBlocking          string = "blocked_due_to_attempt_threshold"
	MultipleUniqueCardsBlocking       string = "blocked_due_to_multiple_unique_cards"
	ConductedExcessiveLoadingBlocking string = "blocked_due_to_conducted_excessive_loading"
)
