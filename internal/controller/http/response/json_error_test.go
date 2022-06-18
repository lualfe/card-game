package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJSONError(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		statusCode int
		wantBody   Error
	}{
		{
			name:       "Bad Request",
			message:    "bad",
			statusCode: http.StatusBadRequest,
			wantBody: Error{
				Message: "bad",
			},
		},
		{
			name:       "Internal Server Error",
			message:    "server error",
			statusCode: http.StatusInternalServerError,
			wantBody: Error{
				Message: "server error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JSONError(w, tt.message, tt.statusCode)

			resp := w.Result()
			defer resp.Body.Close()

			code := resp.StatusCode
			if code != tt.statusCode {
				t.Fatalf("JSONError() | got status code %d, want %d", code, tt.statusCode)
			}

			var got Error
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.wantBody); diff != "" {
				t.Fatalf("JSONError() | (-got +want):\n%s", diff)
			}
		})
	}
}
