package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Timeout 관리하는 미들웨어
func TestTimeoutMiddleware(t *testing.T) {
	handler := http.TimeoutHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
			// 의도적으로 타임아웃 만들기
			time.Sleep(time.Minute)
		}),
		time.Second,
		"Timed out while reading response",
	)

	r := httptest.NewRequest(http.MethodGet, "http://test/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("unexpected status code: %q", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if actual := string(b); actual != "Timed out while reading response" {
		t.Logf("unexpected body: %q", actual)
	}
}
