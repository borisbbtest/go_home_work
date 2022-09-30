package handlers

import (
	"net/http"
	"testing"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/v5"
)

func TestWrapperHandler_FileServer(t *testing.T) {
	type fields struct {
		ServerConf *config.ServiceShortURLConfig
		Storage    storage.Storage
		UserID     string
	}
	type args struct {
		r    chi.Router
		path string
		root http.FileSystem
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
			hook := &WrapperHandler{
				ServerConf: tt.fields.ServerConf,
				Storage:    tt.fields.Storage,
				UserID:     tt.fields.UserID,
			}
			hook.FileServer(tt.args.r, tt.args.path, tt.args.root)
		})
	}
}
