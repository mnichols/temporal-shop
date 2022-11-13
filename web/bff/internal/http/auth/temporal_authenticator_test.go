package auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/temporal-shop/services/go/pkg/shopping"
	"go.temporal.io/api/workflowservice/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authenticatorTestCase struct {
	desc             string
	path             string
	expectBody       string
	expectStatusCode int
	email            string
	sessionToken     string
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

func Test_Temporal_AuthenticateRequest(t *testing.T) {
	const encKey = "feefifum"
	const email = "mike@example.com"
	sessionToken, err := shopping.GenerateShopperHash(encKey, email)
	if err != nil {
		t.Fatal(err)
	}

	cases := []authenticatorTestCase{
		{
			desc:             "valid session token is present",
			path:             "/",
			expectStatusCode: http.StatusOK,
			expectBody:       http.StatusText(http.StatusOK),
			sessionToken:     sessionToken,
			email:            email,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.desc, func(t *testing.T) {
			A := assert.New(t)
			s := &mockTemporal{}
			s.On("DescribeWorkflowExecution", mock.AnythingOfType("context.Context"), testCase.email, "").Return(&workflowservice.DescribeWorkflowExecutionResponse{}, nil)
			session := NewTemporalSessionStore(s)
			r := httptest.NewRequest(http.MethodGet, testCase.path, nil)
			if testCase.sessionToken != "" {
				r.AddCookie(&http.Cookie{Name: sessionCookieName, Value: testCase.sessionToken})
			}
			sut, err := NewAuthenticator(encKey, session)
			A.NoError(err)
			auth, err := sut.AuthenticateRequest(r)
			s.AssertExpectations(t)
			A.NoError(err)
			A.Equal(testCase.email, auth.Email)

		})
	}
	//router := chi.NewRouter()
	//testserver := httptest.NewServer(router)
	//router.With(Authenticate(encKey)).Handle()
	//defer testserver.Close()
	//u, err := url.Parse(testserver.URL)
	//if err != nil {
	//	t.Fatal("unable to server url")
	//}
	//for _, testCase := range cases {
	//	t.Run(testCase.desc, func(t *testing.T) {
	//		httpClient := &http.Client{
	//			Transport: nil,
	//			CheckRedirect: func(req *http.Request, via []*http.Request) error {
	//				return http.ErrUseLastResponse
	//			},
	//			Jar:     nil,
	//			Timeout: 0,
	//		}
	//		p, err := url.Parse(testCase.path)
	//		if err != nil {
	//			t.Fatal("unable to parse path")
	//		}
	//		u = u.ResolveReference(p)
	//		resp, err := httpClient.Get(u.String())
	//		if err != nil {
	//			t.Fatalf("failed to GET: %v", err)
	//		}
	//		if resp.StatusCode != testCase.expectStatusCode {
	//			t.Errorf("handler returned wrong status code: got %v want %v",
	//				resp.StatusCode, testCase.expectStatusCode)
	//		}
	//		bytes, err := io.ReadAll(resp.Body)
	//		if err != nil {
	//			t.Fatalf("failed to read body %v", err)
	//		}
	//		if testCase.expectBody != "" {
	//			if string(bytes) != testCase.expectBody {
	//				t.Errorf("handler returned wrong body: got %v want %v", string(bytes), testCase.expectBody)
	//			}
	//		}
	//	})
	//}

}
