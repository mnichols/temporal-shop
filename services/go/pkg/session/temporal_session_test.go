package session

import (
	"context"
	"github.com/stretchr/testify/mock"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/mocks"
)

type temporalSessionTestCase struct {
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
func (m mockTemporal) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, p ...interface{}) (client.WorkflowRun, error) {
	args := m.Called(ctx, options, workflow, p)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return &mocks.WorkflowRun{}, nil
}

//func Test_Temporal_AuthenticateRequest(t *testing.T) {
//	const encKey = "feefifum"
//	const email = "mike@example.com"
//	var associatedData = encrypt.AssociatedData
//	sessionToken, err := shopping.GenerateShopperHash(encKey, email, associatedData)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	cases := []temporalSessionTestCase{
//		{
//			desc:         "valid session cookie token is present but shopper is not found",
//			path:         "/",
//			sessionToken: &sessionToken,
//			email:        email,
//			temporalErr:  &serviceerror.NotFound{},
//			expectErr:    auth.AuthenticationFailedError,
//		},
//		{
//			desc:         "valid session cookie token is present and shopper is found",
//			path:         "/",
//			sessionToken: &sessionToken,
//			email:        email,
//			temporalErr:  nil,
//			expectErr:    nil,
//		},
//		{
//			desc:         "valid session cookie token is not present",
//			path:         "/",
//			sessionToken: nil,
//			email:        email,
//			expectErr:    auth.AuthenticationFailedError,
//		},
//		{
//			desc:         "valid auth header token is present but shopper is not found",
//			path:         "/",
//			sessionToken: &sessionToken,
//			email:        email,
//			temporalErr:  &serviceerror.NotFound{},
//			expectErr:    auth.AuthenticationFailedError,
//			header:       true,
//		},
//		{
//			desc:         "valid auth header token is present and shopper is found",
//			path:         "/",
//			sessionToken: &sessionToken,
//			email:        email,
//			temporalErr:  nil,
//			expectErr:    nil,
//			header:       true,
//		},
//		{
//			desc:         "valid auth header token is not present",
//			path:         "/",
//			sessionToken: nil,
//			email:        email,
//			expectErr:    auth.AuthenticationFailedError,
//			header:       true,
//		},
//	}
//	for _, testCase := range cases {
//		t.Run(testCase.desc, func(t *testing.T) {
//			A := assert.New(t)
//			id, ierr := session2.NewID([]byte(encKey), associatedData, &values.SessionID{Email: testCase.email})
//			A.NoError(ierr)
//			//s := &mockTemporal{}
//			s := &mocks.Client{}
//			session := NewTemporalSessionStore(s)
//			r := httptest.NewRequest(http.MethodGet, testCase.path, nil)
//			if testCase.sessionToken != nil {
//				cmd := &commands.RefreshShopperRequest{
//					LastSeenAt: nil,
//					Email:      email,
//				}
//				s.On(
//					"SignalWorkflow",
//					mock.Anything,
//					id.ShopperID(),
//					"",
//					orchestrations.SignalName(cmd),
//					mock.MatchedBy(func(in *commands.RefreshShopperRequest) bool {
//						return in.Email == email
//					}),
//				).Return(testCase.temporalErr)
//				//if testCase.header {
//				//	A.NoError(auth.tokenizeRequest(encKey, email, r))
//				//} else {
//				//	r.AddCookie(&http.Cookie{Name: auth.sessionCookieName, Value: *testCase.sessionToken})
//				//}
//			}
//			A.NoError(session.Validate(context.Background(), email))
//
//			//sut, err := auth.NewAuthenticator(encKey, associatedData, session)
//			//A.NoError(err)
//			//auth, err := sut.AuthenticateRequest(r)
//			//s.AssertExpectations(t)
//			//if testCase.expectErr != nil {
//			//	A.EqualError(err, testCase.expectErr.Error())
//			//} else {
//			//	A.NoError(err)
//			//	A.Equal(testCase.email, auth.Email)
//			//}
//			//s.AssertExpectations(t)
//		})
//	}
//
//}
