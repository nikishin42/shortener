package server

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"

	"github.com/nikishin42/shortener/cmd/shortener/pkg/shortener"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/storage"
)

func TestServer_Homepage(t *testing.T) {
	type fields struct {
		Cache     storage.StorageI
		Shortener shortener.ShortenerI
		Router    *mux.Router
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type expext struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expext
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New()
			a.Homepage(tt.args.w, tt.args.r)
		})
	}
}

func TestServer_Redirect(t *testing.T) {
	type fields struct {
		Cache     storage.StorageI
		Shortener shortener.ShortenerI
		Router    *mux.Router
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Server{
				Storage:   tt.fields.Cache,
				Shortener: tt.fields.Shortener,
				Router:    tt.fields.Router,
			}
			a.Redirect(tt.args.w, tt.args.r)
		})
	}
}
