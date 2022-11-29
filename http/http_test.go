package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetText(t *testing.T) {

	tests := []struct {
		serverStatus   int
		serverResponse []byte
		want           []byte
	}{
		{http.StatusOK, []byte("<p>Hello World</p>"), []byte("Hello World\n")},
		//{http.StatusOK, []byte("<p>Hello World</p>"), []byte("wrong")},
		{http.StatusOK, nil, []byte("\n")},
	}

	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.serverStatus)
			w.Write(tc.serverResponse)
		}))
		got, err := GetText(server.URL)
		if err != nil {
			t.Errorf("The error is %v", err)
		} else if diff := cmp.Diff(string(tc.want), string(got)); diff != "" {
			t.Errorf("diff %s", diff)
		}

		server.Close()
	}
}
