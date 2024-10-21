package router

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Request struct {
	method string
	path   string
	body   io.Reader
}

type Test struct {
	name    string
	request Request
	route   Route
	want    bool
}

func TestRouter_resolveRoute(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {}

	tests := []Test{
		{
			name:    "root path positive test",
			request: Request{method: http.MethodGet, path: "/", body: nil},
			route:   Route{Method: http.MethodGet, Path: "/", Handler: mockHandler},
			want:    true,
		},
		{
			name:    "mismatch route negative test",
			request: Request{method: http.MethodGet, path: "/user", body: nil},
			route:   Route{Method: http.MethodGet, Path: "/profile", Handler: mockHandler},
			want:    false,
		},
		{
			name:    "regexp match positive test",
			request: Request{method: http.MethodGet, path: "/user/1", body: nil},
			route:   Route{Method: http.MethodGet, Path: `/user/\d+`, Handler: mockHandler},
			want:    true,
		},
		{
			name:    "regexp mismatch negative test",
			request: Request{method: http.MethodGet, path: "/user/1", body: nil},
			route:   Route{Method: http.MethodGet, Path: `/profile/\d+`, Handler: mockHandler},
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.request.method, test.request.path, test.request.body)

			assert.Equal(t, test.want, NewRouter().resolveRoute(test.route, req))
		})
	}
}
