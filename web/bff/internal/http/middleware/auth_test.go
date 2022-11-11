package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type middlewareTestCase struct {
	desc             string
	path             string
	expectBody       string
	expectStatusCode int
	next             *next
}
type next struct {
	invoked bool
}

func (n *next) handle(w http.ResponseWriter, r *http.Request) {
	n.invoked = true
	w.WriteHeader(http.StatusOK)
}
func TestAuthenticationRequired(t *testing.T) {
	const encKey = "feefifum"

	cases := []middlewareTestCase{
		{
			desc:             "valid session token is present",
			path:             "/",
			expectStatusCode: http.StatusUnauthorized,
			expectBody:       http.StatusText(http.StatusUnauthorized),
			next:             &next{},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.desc, func(t *testing.T) {

			sut := Authenticate(encKey)(http.HandlerFunc(testCase.next.handle))
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, testCase.path, nil)
			sut.ServeHTTP(w, r)
			A := assert.New(t)
			A.Equal(testCase.expectStatusCode, w.Code)
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
