package inmemory_test

import (
	"context"
	"testing"
	"wallet/internal"
	"wallet/internal/inmemory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionRepository(t *testing.T) {
	got := inmemory.NewTransactionRepository()
	assert.NotNil(t, got)
}

var _ = Describe("transaction repository", func() {
	repo := inmemory.NewTransactionRepository()

	Context("store debit transaction data", func() {
		data := internal.TransactionData{
			UserID: "test",
			Amount: 10000,
			Type:   internal.DebitTransaction,
		}
		It("will return no error", func() {
			err := repo.Store(context.Background(), data)
			Expect(err).Should(BeNil())
		})

		It("will found in FindTopTransactionsByUserID", func() {
			res := repo.FindTopTransactionsByUserID(context.Background(), data.UserID, 1)
			Expect(res).Should(HaveLen(1))
			Expect(res[0].UserID).Should(Equal(data.UserID))
			Expect(res[0].Amount).Should(Equal(data.Amount))
			Expect(res[0].Type).Should(Equal(data.Type))
		})

		It("will found in FindTopTransactingUser", func() {
			res := repo.FindTopTransactingUser(context.Background(), 1)
			Expect(res).Should(HaveLen(1))
			Expect(res[0].UserID).Should(Equal(data.UserID))
			Expect(res[0].Amount).Should(Equal(data.Amount))
			Expect(res[0].Type).Should(Equal(data.Type))
		})
	})

	Context("store credit transaction data", func() {
		data := internal.TransactionData{
			UserID: "test",
			Amount: 100,
			Type:   internal.CreditTransaction,
		}
		It("will return no error", func() {
			err := repo.Store(context.Background(), data)
			Expect(err).Should(BeNil())
		})

		It("will found in last index of FindTopTransactionsByUserID result", func() {
			res := repo.FindTopTransactionsByUserID(context.Background(), data.UserID, 2)
			Expect(res).Should(HaveLen(2))

			Expect(res[len(res)-1]).Should(MatchFields(IgnoreExtras, Fields{
				"UserID": Equal(data.UserID),
				"Amount": Equal(data.Amount),
				"Type":   Equal(data.Type),
			}))
		})

		It("will not found in FindTopTransactingUser", func() {
			res := repo.FindTopTransactingUser(context.Background(), 1)
			Expect(res).Should(HaveLen(1))
			Expect(res).ShouldNot(ContainElement(MatchFields(IgnoreExtras, Fields{
				"UserID": Equal(data.UserID),
				"Amount": Equal(data.Amount),
				"Type":   Equal(data.Type),
			})))
		})
	})
})
