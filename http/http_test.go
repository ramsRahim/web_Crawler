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

func TestReadRss(t *testing.T) {

	tests := []struct {
		serverStatus   int
		serverResponse []byte
		want           []string
		wantErr        bool
	}{
		{http.StatusOK, []byte("<item><link>Hello World</link></item><item><link>Rahim</link></item>"), []string{"Hello World", "Rahim"}, false},
		{http.StatusOK, []byte("<item><link>Hello World</link</item><item><link>Rahim</link</item>"), nil, true},
	}

	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.serverStatus)
			w.Write(tc.serverResponse)
		}))
		got, err := ReadRSS(server.URL)
		if !tc.wantErr && err != nil {
			t.Errorf("didn't expect error but returned one, %v", err)
		} else if len(got) != len(tc.want) {
			t.Errorf("output didn't match")
		} else {
			for i, g := range tc.want {
				if diff := cmp.Diff(g, got[i]); diff != "" {
					t.Errorf("diff %s", diff)
				}
			}
		}

		server.Close()
	}
}

func TestError(t *testing.T) {

	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://www.cobaltspeech.com", false},
		{"", true},
		//{"https://jsonplaceholder.typicode.com/posts", true},
	}

	for _, tc := range tests {
		_, err1 := GetText(tc.url)
		_, err2 := ReadRSS(tc.url)

		if !tc.wantErr && err1 != nil && err2 != nil {
			t.Errorf("didn't expect an error but returned one")
		} else if tc.wantErr && err1 == nil && err2 == nil {
			t.Errorf("expected an error but didn't return one")
		}
	}
}
