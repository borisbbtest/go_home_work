package handlers_test

import (
	"net/http"
	"testing"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/storage"
)

func TestWrapperHandler_GetHandler(t *testing.T) {
	type fields struct {
		URLStore   storage.StoreDB
		ServerConf *config.ServiceShortURLConfig
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
			hook := &WrapperHandler{
				URLStore:   tt.fields.URLStore,
				ServerConf: tt.fields.ServerConf,
			}
			hook.GetHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestWrapperHandler_PostHandler(t *testing.T) {
	type fields struct {
		URLStore   storage.StoreDB
		ServerConf *config.ServiceShortURLConfig
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
			hook := &WrapperHandler{
				URLStore:   tt.fields.URLStore,
				ServerConf: tt.fields.ServerConf,
			}
			hook.PostHandler(tt.args.w, tt.args.r)
		})
	}
}
