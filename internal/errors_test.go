package internal_test

import (
	"errors"
	"fmt"
	"testing"
	"wallet/internal"

	"github.com/stretchr/testify/assert"
)

func TestWrapErr(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				err: errors.New("another error"),
				msg: "random",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := internal.WrapErr(tt.args.err, tt.args.msg)
			if !errors.Is(got, tt.args.err) {
				t.Fail()
				return
			}

			assert.Equal(t, fmt.Sprintf("%s. reason: %s", tt.args.err, tt.args.msg), got.Error())
		})
	}
}
