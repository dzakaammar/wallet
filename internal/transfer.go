package internal

type P2PTransferRequest struct {
	InitiatorUserID string
	ToUsername      string
	Amount          int
}

func (p P2PTransferRequest) Validate() error {
	if p.InitiatorUserID == "" {
		return WrapErr(ErrInvalidParameter, "invalid user")
	}

	if p.ToUsername == "" {
		return WrapErr(ErrInvalidParameter, "invalid to username param")
	}

	if p.Amount <= 0 {
		return WrapErr(ErrInvalidParameter, "invalid amount")
	}

	return nil
}

type TopUpRequest struct {
	UserID string
	Amount int
}

func (t TopUpRequest) Validate() error {
	if t.UserID == "" {
		return WrapErr(ErrInvalidParameter, "invalid user")
	}

	if t.Amount <= 0 || t.Amount >= 10000000 {
		return WrapErr(ErrInvalidParameter, "invalid amount")
	}

	return nil
}

type TransferEvent struct {
	DebitUserID  string
	CreditUserID string
	Amount       int
}

func (t TransferEvent) ToTransactionData() []TransactionData {
	return []TransactionData{
		{
			UserID: t.DebitUserID,
			Amount: t.Amount,
			Type:   DebitTransaction,
		},
		{
			UserID: t.CreditUserID,
			Amount: t.Amount,
			Type:   CreditTransaction,
		},
	}
}
