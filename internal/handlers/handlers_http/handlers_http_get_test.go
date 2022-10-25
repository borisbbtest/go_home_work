package handlershttp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/borisbbtest/go_home_work/internal/config"
	handlershttp "github.com/borisbbtest/go_home_work/internal/handlers/handlers_http"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func TestWrapperHandler_GetHandlerCooke(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Connetion test TestDeleteURLHandlers",
			want: want{
				code:        204,
				response:    `No Content`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

			th := handlershttp.WrapperHandler{
				ServerConf: &config.ServiceShortURLConfig{
					Port:       8080,
					ServerHost: "localhost",
					BaseURL:    "http://localhost:8080",
				},
			}
			th.Storage, _ = storage.NewFileStorage("")

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			r := chi.NewRouter()
			r.Use(middleware.RequestID)
			r.Use(middleware.RealIP)
			r.Use(middleware.Logger)
			r.Use(th.GzipHandle)
			r.Use(th.MidSetCookie)
			r.Get("/api/user/urls", th.GetHandlerCooke)

			r.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

		})
	}
}

func TestWrapperHandler_GetHandlerPing(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Connetion test TestDeleteURLHandlers",
			want: want{
				code:        500,
				response:    `error connection`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/ping", nil)

			th := handlershttp.WrapperHandler{
				ServerConf: &config.ServiceShortURLConfig{
					Port:       8080,
					ServerHost: "localhost",
					BaseURL:    "http://localhost:8080",
				},
			}
			th.Storage, _ = storage.NewFileStorage("")

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(th.GetHandlerPing)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// // заголовок ответа
			// if res.Header.Get("Content-Type") != tt.want.contentType {
			// 	//	t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			// }
		})
	}
}
