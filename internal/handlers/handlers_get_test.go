package handlers

import (
	"net/http"
	"testing"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/storage"
)

func TestWrapperHandler_GetHandler(t *testing.T) {
	type fields struct {
		ServerConf *config.ServiceShortURLConfig
		Storage    storage.Storage
		UserID     string
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
				ServerConf: tt.fields.ServerConf,
				Storage:    tt.fields.Storage,
				UserID:     tt.fields.UserID,
			}
			hook.GetHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestWrapperHandler_GetHandlerCooke(t *testing.T) {
	type fields struct {
		ServerConf *config.ServiceShortURLConfig
		Storage    storage.Storage
		UserID     string
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
				ServerConf: tt.fields.ServerConf,
				Storage:    tt.fields.Storage,
				UserID:     tt.fields.UserID,
			}
			hook.GetHandlerCooke(tt.args.w, tt.args.r)
		})
	}
}

func TestWrapperHandler_GetHandlerPing(t *testing.T) {
	type fields struct {
		ServerConf *config.ServiceShortURLConfig
		Storage    storage.Storage
		UserID     string
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
				ServerConf: tt.fields.ServerConf,
				Storage:    tt.fields.Storage,
				UserID:     tt.fields.UserID,
			}
			hook.GetHandlerPing(tt.args.w, tt.args.r)
		})
	}
}
