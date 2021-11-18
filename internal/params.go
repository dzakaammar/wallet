package internal

type TransactionPerUserResponse struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}

type UserTransactionsResponse struct {
	Username        string `json:"username"`
	TransactedValue int    `json:"transacted_value"`
}

type TransferRequest struct {
	InitiatorUserID string
	ToUsername      string
	Amount          int
}
