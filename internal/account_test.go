package internal_test

import (
	"errors"
	"testing"
	"wallet/internal"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    *internal.Account
		wantErr error
	}{
		{
			name: "OK",
			args: args{
				userID: "test",
			},
			want: &internal.Account{
				UserID:  "test",
				Balance: 0,
			},
		},
		{
			name: "Not OK",
			args: args{
				userID: "",
			},
			want:    nil,
			wantErr: internal.ErrInvalidParameter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := internal.NewAccount(tt.args.userID)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fail()
				}
				assert.Nil(t, tt.want)
				return
			}
			assert.NotEmpty(t, got.ID)
			assert.Equal(t, tt.want.UserID, got.UserID)
			assert.Equal(t, tt.want.Balance, got.Balance)
		})
	}
}

func TestAccount_TransferTo(t *testing.T) {
	type fields struct {
		ID      string
		UserID  string
		Balance int
	}
	type args struct {
		creditur *internal.Account
		amount   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "OK",
			fields: fields{
				ID:      "debiturAccountID",
				UserID:  "debiturID",
				Balance: 1001,
			},
			args: args{
				creditur: &internal.Account{
					ID:      "crediturAccountID",
					UserID:  "crediturID",
					Balance: 0,
				},
				amount: 1000,
			},
		},
		{
			name: "Not OK: insufficient balance",
			fields: fields{
				ID:      "debiturAccountID",
				UserID:  "debiturID",
				Balance: 1000,
			},
			args: args{
				creditur: &internal.Account{
					ID:      "crediturAccountID",
					UserID:  "crediturID",
					Balance: 0,
				},
				amount: 1000,
			},
			wantErr: internal.ErrInsufficientBalance,
		},
		{
			name: "Not OK: same id",
			fields: fields{
				ID:      "debiturAccountID",
				UserID:  "debiturID",
				Balance: 1000,
			},
			args: args{
				creditur: &internal.Account{
					ID:      "debiturAccountID",
					UserID:  "debiturID",
					Balance: 0,
				},
				amount: 1000,
			},
			wantErr: internal.ErrInvalidParameter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &internal.Account{
				ID:      tt.fields.ID,
				UserID:  tt.fields.UserID,
				Balance: tt.fields.Balance,
			}

			crediturPrevBalance := tt.args.creditur.Balance

			err := a.TransferTo(tt.args.creditur, tt.args.amount)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fail()
				}
				return
			}

			assert.Equal(t, tt.fields.Balance-tt.args.amount, a.Balance)
			assert.Equal(t, crediturPrevBalance+tt.args.amount, tt.args.creditur.Balance)
		})
	}
}

func TestAccount_Deposit(t *testing.T) {
	type fields struct {
		ID      string
		UserID  string
		Balance int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "OK",
			fields: fields{
				ID:      "accountID",
				UserID:  "userID",
				Balance: 0,
			},
			args: args{
				amount: 1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &internal.Account{
				ID:      tt.fields.ID,
				UserID:  tt.fields.UserID,
				Balance: tt.fields.Balance,
			}
			a.Deposit(tt.args.amount)
			assert.Equal(t, tt.fields.Balance+tt.args.amount, a.Balance)
		})
	}
}
