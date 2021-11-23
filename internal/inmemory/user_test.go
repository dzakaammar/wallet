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

func TestNewUserRepository(t *testing.T) {
	got := inmemory.NewUserRepository()
	assert.NotNil(t, got)
}

var _ = Describe("user repository", func() {
	repo := inmemory.NewUserRepository()

	data := internal.User{
		ID:       "test",
		Username: "username",
	}
	Context("store user data", func() {
		It("will return no error", func() {
			err := repo.Store(context.Background(), data)
			Expect(err).Should(BeNil())
		})
	})

	Context("store duplicate user data", func() {
		It("will return error", func() {
			err := repo.Store(context.Background(), data)
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(MatchError(internal.ErrDataAlreadyExists))
		})
	})

	Context("find by id", func() {
		It("will return the data", func() {
			res, err := repo.FindByID(context.Background(), data.ID)
			Expect(err).Should(BeNil())
			Expect(*res).Should(Equal(data))
		})
	})

	Context("find by username", func() {
		It("will return the data", func() {
			res, err := repo.FindByUsername(context.Background(), data.Username)
			Expect(err).Should(BeNil())
			Expect(*res).Should(Equal(data))
		})
	})

	Context("find not exists data by id", func() {
		It("will return error", func() {
			res, err := repo.FindByID(context.Background(), "anotherID")
			Expect(res).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(MatchError(internal.ErrDataNotFound))
		})
	})

	Context("find not exists data by username", func() {
		It("will return error", func() {
			res, err := repo.FindByUsername(context.Background(), "anotherUsername")
			Expect(res).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(MatchError(internal.ErrDataNotFound))
		})
	})
})
