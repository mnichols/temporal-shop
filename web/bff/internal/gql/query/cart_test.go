package query

import (
	"github.com/stretchr/testify/assert"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/format"
	"testing"
)

func Test_centsToDollars(t *testing.T) {
	type args struct {
		val int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "converts to dollars ok",
		args: args{val: 198763},
		want: "1987.63",
	},
		{
			name: "pads when less than a dollar",
			args: args{val: 50},
			want: "0.50",
		},
		{
			name: "supports single digits",
			args: args{val: 5},
			want: "0.05",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, format.CentsToDollars(tt.args.val), "centsToDollars(%v)", tt.args.val)
		})
	}
}
func Test_bpsToPercent(t *testing.T) {
	type args struct {
		val int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "converts to percent ok",
		args: args{val: 425},
		want: "4.25",
	},
		{
			name: "doesnt pad when less than 100",
			args: args{val: 50},
			want: "50",
		},
		{
			name: "supports single digits",
			args: args{val: 5},
			want: "5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, format.BpsToPercent(tt.args.val), "bpsToPercent(%v)", tt.args.val)
		})
	}
}
