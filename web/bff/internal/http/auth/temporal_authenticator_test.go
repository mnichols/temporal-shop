package auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/temporal-shop/api/temporal_shop/commands/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/orchestrations"
	"github.com/temporalio/temporal-shop/services/go/pkg/shopping"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authenticatorTestCase struct {
	desc         string
	path         string
	email        string
	sessionToken *string
	expectErr    error
	temporalErr  error
	header       bool
}
type mockTemporal struct {
	mock.Mock
}

func (m mockTemporal) DescribeWorkflowExecution(ctx context.Context, s string, s2 string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
	args := m.Called(ctx, s, s2)
	out := args.Get(0)
	if out != nil {
		return out.(*workflowservice.DescribeWorkflowExecutionResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m mockTemporal) SignalWorkflow(ctx context.Context, wid, rid, signalName string, arg interface{}) error {
	args := m.Called(ctx, wid, rid, signalName, arg)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func Test_Temporal_AuthenticateRequest(t *testing.T) {
	const encKey = "feefifum"
	const email = "mike@example.com"
	sessionToken, err := shopping.GenerateShopperHash(encKey, email)
	if err != nil {
		t.Fatal(err)
	}

	cases := []authenticatorTestCase{
		{
			desc:         "valid session cookie token is present but shopper is not found",
			path:         "/",
			sessionToken: &sessionToken,
			email:        email,
			temporalErr:  &serviceerror.NotFound{},
			expectErr:    AuthenticationFailedError,
		},
		{
			desc:         "valid session cookie token is present and shopper is found",
			path:         "/",
			sessionToken: &sessionToken,
			email:        email,
			temporalErr:  nil,
			expectErr:    nil,
		},
		{
			desc:         "valid session cookie token is not present",
			path:         "/",
			sessionToken: nil,
			email:        email,
			expectErr:    AuthenticationFailedError,
		},
		{
			desc:         "valid auth header token is present but shopper is not found",
			path:         "/",
			sessionToken: &sessionToken,
			email:        email,
			temporalErr:  &serviceerror.NotFound{},
			expectErr:    AuthenticationFailedError,
			header:       true,
		},
		{
			desc:         "valid auth header token is present and shopper is found",
			path:         "/",
			sessionToken: &sessionToken,
			email:        email,
			temporalErr:  nil,
			expectErr:    nil,
			header:       true,
		},
		{
			desc:         "valid auth header token is not present",
			path:         "/",
			sessionToken: nil,
			email:        email,
			expectErr:    AuthenticationFailedError,
			header:       true,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.desc, func(t *testing.T) {
			A := assert.New(t)
			s := &mockTemporal{}
			session := NewTemporalSessionStore(s)
			r := httptest.NewRequest(http.MethodGet, testCase.path, nil)
			if testCase.sessionToken != nil {
				cmd := &commands.RefreshShopperRequest{
					LastSeenAt: nil,
					Email:      email,
				}
				s.On(
					"SignalWorkflow",
					mock.Anything,
					testCase.email,
					"",
					orchestrations.SignalName(cmd),
					mock.MatchedBy(func(in *commands.RefreshShopperRequest) bool {
						return in.Email == email
					}),
				).Return(testCase.temporalErr)
				if testCase.header {
					A.NoError(tokenizeRequest(encKey, email, r))
				} else {
					r.AddCookie(&http.Cookie{Name: sessionCookieName, Value: *testCase.sessionToken})
				}
			}
			sut, err := NewAuthenticator(encKey, session)
			A.NoError(err)
			auth, err := sut.AuthenticateRequest(r)
			s.AssertExpectations(t)
			if testCase.expectErr != nil {
				A.EqualError(err, testCase.expectErr.Error())
			} else {
				A.NoError(err)
				A.Equal(testCase.email, auth.Email)
			}
			s.AssertExpectations(t)
		})
	}

}
