package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJSON(t *testing.T) {
	type fakeResponse struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name       string
		obj        fakeResponse
		statusCode int
	}{
		{
			name:       "Status OK",
			statusCode: http.StatusOK,
			obj: fakeResponse{
				Message: "ok",
			},
		},
		{
			name:       "Status Created",
			statusCode: http.StatusAccepted,
			obj: fakeResponse{
				Message: "created",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			JSON(w, tt.obj, tt.statusCode)

			resp := w.Result()
			defer resp.Body.Close()

			code := resp.StatusCode
			if code != tt.statusCode {
				t.Fatalf("JSON() | got status code %d, want %d", code, tt.statusCode)
			}

			var got fakeResponse
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if err = json.Unmarshal(body, &got); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.obj); diff != "" {
				t.Fatalf("JSON() | (-got +want):\n%s", diff)
			}
		})
	}
}
