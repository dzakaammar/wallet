package internal_test

import (
	"errors"
	"reflect"
	"testing"
	"wallet/internal"
)

func TestP2PTransferRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       internal.P2PTransferRequest
		wantErr error
	}{
		{
			name: "OK",
			p: internal.P2PTransferRequest{
				InitiatorUserID: "test",
				ToUsername:      "username",
				Amount:          1000,
			},
		},
		{
			name: "Not OK: invalid initiator user id",
			p: internal.P2PTransferRequest{
				InitiatorUserID: "",
				ToUsername:      "username",
				Amount:          1000,
			},
			wantErr: internal.ErrInvalidParameter,
		},
		{
			name: "Not OK: invalid to username",
			p: internal.P2PTransferRequest{
				InitiatorUserID: "test",
				ToUsername:      "",
				Amount:          1000,
			},
			wantErr: internal.ErrInvalidParameter,
		},
		{
			name: "Not OK: invalid amount",
			p: internal.P2PTransferRequest{
				InitiatorUserID: "test",
				ToUsername:      "username",
				Amount:          0,
			},
			wantErr: internal.ErrInvalidParameter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fail()
			}
		})
	}
}

func TestTopUpRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tr      internal.TopUpRequest
		wantErr error
	}{
		{
			name: "OK",
			tr: internal.TopUpRequest{
				UserID: "userID",
				Amount: 1000,
			},
		},
		{
			name: "Not OK: invalid userID",
			tr: internal.TopUpRequest{
				UserID: "",
				Amount: 1000,
			},
			wantErr: internal.ErrInvalidParameter,
		},
		{
			name: "Not OK: amount gt 10000000",
			tr: internal.TopUpRequest{
				UserID: "userID",
				Amount: 10000001,
			},
			wantErr: internal.ErrInvalidParameter,
		},
		{
			name: "Not OK: amount eq 10000000",
			tr: internal.TopUpRequest{
				UserID: "userID",
				Amount: 10000000,
			},
			wantErr: internal.ErrInvalidParameter,
		},
		{
			name: "Not OK: amount equal 0",
			tr: internal.TopUpRequest{
				UserID: "userID",
				Amount: 0,
			},
			wantErr: internal.ErrInvalidParameter,
		},
		{
			name: "Not OK: amount lt 0",
			tr: internal.TopUpRequest{
				UserID: "userID",
				Amount: -1,
			},
			wantErr: internal.ErrInvalidParameter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tr.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fail()
			}
		})
	}
}

func TestTransferEvent_ToTransactionData(t *testing.T) {
	tests := []struct {
		name string
		tr   internal.TransferEvent
		want []internal.TransactionData
	}{
		{
			name: "OK",
			tr: internal.TransferEvent{
				DebitUserID:  "debitUserID",
				CreditUserID: "creditUserID",
				Amount:       1000,
			},
			want: []internal.TransactionData{
				{
					UserID: "debitUserID",
					Amount: 1000,
					Type:   internal.DebitTransaction,
				},
				{
					UserID: "creditUserID",
					Amount: 1000,
					Type:   internal.CreditTransaction,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToTransactionData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransferEvent.ToTransactionData() = %v, want %v", got, tt.want)
			}
		})
	}
}
