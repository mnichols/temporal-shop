package orchestrations

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SignalName(t *testing.T) {
	A := assert.New(t)
	expect := "temporal_shop.test.v1.SignalRequest"
	actual := SignalName(&test.SignalRequest{})
	A.Equal(expect, actual)
}
func Test_QueryName(t *testing.T) {
	A := assert.New(t)
	expect := "temporal_shop.test.v1.QueryRequest"
	actual := QueryName(&test.QueryRequest{})
	A.Equal(expect, actual)
}
