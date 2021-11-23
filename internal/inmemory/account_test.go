package inmemory_test

import (
	"context"
	"testing"
	"wallet/internal"
	"wallet/internal/inmemory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountRepository(t *testing.T) {
	got := inmemory.NewAccountRepository()
	assert.NotNil(t, got)
}

var _ = Describe("account repo test", func() {
	repo := inmemory.NewAccountRepository()
	Context("storing data", func() {
		account := internal.Account{
			ID:      "test",
			UserID:  "userID",
			Balance: 100,
		}
		It("will return no error", func() {
			err := repo.Store(context.Background(), account)
			Expect(err).Should(BeNil())
		})

		It("will searchable", func() {
			res, err := repo.FindByUserID(context.Background(), account.UserID)
			Expect(err).Should(BeNil())
			Expect(*res).To(Equal(account))
		})
	})

	Context("storing duplicate data", func() {
		account := internal.Account{
			ID:      "test",
			UserID:  "userID",
			Balance: 100,
		}
		It("will return error", func() {
			err := repo.Store(context.Background(), account)
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(MatchError(internal.ErrDataAlreadyExists))
		})
	})

	Context("updating not existing data", func() {
		It("will return error", func() {
			err := repo.Updates(context.Background(), internal.Account{
				ID:      "notfound",
				UserID:  "userID",
				Balance: 1000,
			})
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(MatchError(internal.ErrDataNotFound))
		})
	})

	Context("updating existing data", func() {
		account := internal.Account{
			ID:      "found",
			UserID:  "anotherUserID",
			Balance: 1000,
		}
		updatedBalance := 2000
		_ = repo.Store(context.Background(), account)
		It("will return no error", func() {
			err := repo.Updates(context.Background(), internal.Account{
				ID:      account.ID,
				UserID:  account.UserID,
				Balance: updatedBalance,
			})
			Expect(err).Should(BeNil())
		})

		It("will change the data", func() {
			res, err := repo.FindByUserID(context.Background(), account.UserID)
			Expect(err).Should(BeNil())
			Expect(res.ID).To(Equal(account.ID))
			Expect(res.UserID).To(Equal(account.UserID))
			Expect(res.Balance).To(Equal(updatedBalance))
		})
	})

	Context("updating existing data", func() {
		account := internal.Account{
			ID:      "found",
			UserID:  "anotherUserID",
			Balance: 1000,
		}
		updatedBalance := 2000
		_ = repo.Store(context.Background(), account)
		It("will return no error", func() {
			err := repo.Updates(context.Background(), internal.Account{
				ID:      account.ID,
				UserID:  account.UserID,
				Balance: updatedBalance,
			})
			Expect(err).Should(BeNil())
		})

		It("will change the data", func() {
			res, err := repo.FindByUserID(context.Background(), account.UserID)
			Expect(err).Should(BeNil())
			Expect(res.ID).To(Equal(account.ID))
			Expect(res.UserID).To(Equal(account.UserID))
			Expect(res.Balance).To(Equal(updatedBalance))
		})
	})
})
