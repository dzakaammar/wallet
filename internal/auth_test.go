package internal_test

import (
	"testing"
	"wallet/internal"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthToken(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				secret: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := internal.NewAuthToken(tt.args.secret)
			assert.NotNil(t, got)
		})
	}
}

func TestAuthToken_GetTokenFor_ParseToken(t *testing.T) {
	type fields struct {
		secret string
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "OK",
			fields: fields{
				secret: "test",
			},
			args: args{
				userID: "userID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := internal.NewAuthToken(tt.fields.secret)

			got, err := a.GetTokenFor(tt.args.userID)
			assert.NoError(t, err)

			claims, err := a.ParseToken(got)
			assert.NoError(t, err)

			assert.Equal(t, tt.args.userID, claims.UserID)
		})
	}
}
