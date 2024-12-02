package lowercase

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		url        string
		next       func(t *testing.T) http.Handler
		assertFunc func(t *testing.T, rw *httptest.ResponseRecorder)
	}{
		{
			name:   "base case",
			config: &Config{},
			url:    "http://localhost/",
			assertFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()
				if rr.Result().StatusCode != http.StatusOK {
					t.Fatalf("expected OK, got %d", rr.Result().StatusCode)
				}
			},
			next: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Helper()
					if req.URL.String() != "http://localhost/" {
						t.Fatalf("wanted path %s, got req.URL: %+v", "http://localhost/", req.URL.String())
					}
				})
			},
		},
		{
			name:   "uppercase redirect",
			config: &Config{},
			url:    "http://localhost/Notthis",
			next: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Helper()
					t.Fatal("got called")
				})
			},
			assertFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()
				if rr.Result().StatusCode != http.StatusMovedPermanently {
					t.Fatalf("expected %d, got %d", http.StatusMovedPermanently, rr.Result().StatusCode)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			handler, err := New(ctx, tt.next(t), tt.config, "lowercase")
			if err != nil {
				t.Fatalf("error with new redirect: %+v", err)
			}
			recorder := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tt.url, nil)
			if err != nil {
				t.Fatalf("error with new request: %+v", err)
			}

			handler.ServeHTTP(recorder, req)
			tt.assertFunc(t, recorder)
		})
	}
}
