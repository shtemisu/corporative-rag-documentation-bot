package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEmbed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embed" {
			t.Errorf("path: want /api/embed, but got %s", r.URL.Path)
		}
		req := embedRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Model != "bge-m3" {
			t.Fatalf("model: want bge-m3, but got %s", req.Model)
		}
		if !reflect.DeepEqual(req.Input, []string{"a", "b"}) {
			t.Fatalf(`want {"a","b"}, but got %+v`, req.Input)
		}
		w.Write([]byte(`{"embeddings":[[0.1,0.2],[0.3,0.4]]}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "bge-m3", "")

	got, err := c.Embed(context.Background(), []string{"a", "b"})
	if err != nil {
		t.Fatalf("failed to get embeddings: %v", err)
	}
	if !reflect.DeepEqual(got, [][]float32{{0.1, 0.2}, {0.3, 0.4}}) {
		t.Errorf("want: %+v, but got: %+v", [][]float32{{0.1, 0.2}, {0.3, 0.4}}, got)
	}
}

func TestEmbedTabled(t *testing.T) {
	testCases := []struct {
		name        string
		texts       []string
		handler     http.HandlerFunc
		wantErrText string
	}{
		{
			name:  "model didn't install",
			texts: []string{"a"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"model 'bge-m3' not found, try pulling it first"}`))
			},
			wantErrText: `ollama returned status 404: {"error":"model 'bge-m3' not found, try pulling it first"}`,
		},
		{
			name:  "bad request syntax (400)",
			texts: []string{"a"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // 400 status
				w.Write([]byte(`{"error":"json: cannot unmarshal number into Go struct field EmbeddingRequest.input of type string"}`))
			},
			wantErrText: `ollama returned status 400: {"error":"json: cannot unmarshal number into Go struct field EmbeddingRequest.input of type string"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			srv := httptest.NewServer(tC.handler)
			defer srv.Close()

			c := NewClient(srv.URL, "bge-m3", "gemma2:2b")

			_, err := c.Embed(context.Background(), tC.texts)
			if err.Error() != tC.wantErrText {
				t.Errorf("\nwant error: %q, but got: %q", tC.wantErrText, err.Error())
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		handler     http.HandlerFunc
		wantErrText string
	}{
		{
			name:   "empty prompt",
			prompt: " ",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"model 'gemma2:2b' not found, try pulling it first"}`))
			},
			wantErrText: `ollama returned status 404: {"error":"model 'gemma2:2b' not found, try pulling it first"}`,
		},
		{
			name:   "bad request",
			prompt: "bruh",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest) // 400 status
				w.Write([]byte(`{"error":"json: cannot unmarshal number into Go struct field EmbeddingRequest.input of type string"}`))
			},
			wantErrText: `ollama returned status 400: {"error":"json: cannot unmarshal number into Go struct field EmbeddingRequest.input of type string"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			c := NewClient(srv.URL, "bge-m3", "gemma2:2b")
			_, err := c.Generate(context.Background(), tt.prompt)
			if err.Error() != tt.wantErrText {
				t.Errorf("\nwant error: %q, but got: %q", tt.wantErrText, err.Error())
			}
		})
	}
}
