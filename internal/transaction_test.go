package internal_test

import (
	"sort"
	"testing"
	"wallet/internal"

	"github.com/stretchr/testify/assert"
)

func TestTransactionData_DisplayAmount(t *testing.T) {
	tests := []struct {
		name string
		tr   internal.TransactionData
		want int
	}{
		{
			name: "OK: debit trx",
			tr: internal.TransactionData{
				UserID: "test",
				Amount: 10000,
				Type:   internal.DebitTransaction,
			},
			want: -10000,
		},
		{
			name: "OK: credit trx",
			tr: internal.TransactionData{
				UserID: "test",
				Amount: 10000,
				Type:   internal.CreditTransaction,
			},
			want: 10000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.DisplayAmount(); got != tt.want {
				t.Errorf("TransactionData.DisplayAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionsData_Sort(t *testing.T) {
	tests := []struct {
		name string
		tr   internal.TransactionsData
		want internal.TransactionsData
		desc bool
	}{
		{
			name: "OK: asc",
			tr: internal.TransactionsData{
				{
					UserID: "userID1",
					Amount: 100,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID2",
					Amount: 99,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID3",
					Amount: 98,
					Type:   internal.CreditTransaction,
				},
			},
			want: internal.TransactionsData{
				{
					UserID: "userID3",
					Amount: 98,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID2",
					Amount: 99,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID1",
					Amount: 100,
					Type:   internal.CreditTransaction,
				},
			},
		},
		{
			name: "OK: desc",
			tr: internal.TransactionsData{
				{
					UserID: "userID3",
					Amount: 98,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID2",
					Amount: 99,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID1",
					Amount: 100,
					Type:   internal.CreditTransaction,
				},
			},
			want: internal.TransactionsData{
				{
					UserID: "userID1",
					Amount: 100,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID2",
					Amount: 99,
					Type:   internal.CreditTransaction,
				},
				{
					UserID: "userID3",
					Amount: 98,
					Type:   internal.CreditTransaction,
				},
			},
			desc: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.desc {
				sort.Sort(sort.Reverse(tt.tr))
			} else {
				sort.Sort(tt.tr)
			}

			assert.EqualValues(t, tt.tr, tt.want)
		})
	}
}
