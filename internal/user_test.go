package internal_test

import (
	"errors"
	"testing"
	"wallet/internal"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := internal.NewUser(tt.args.username)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fail()
					return
				}
				assert.Nil(t, got)
				return
			}

			assert.NotEmpty(t, got.ID)
			assert.Equal(t, tt.args.username, got.Username)
		})
	}
}
