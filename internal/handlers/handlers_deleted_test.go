package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/handlers"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func TestWrapperHandler_DeleteURLHandlers(t *testing.T) {
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
			name: "Connetion test #1",
			want: want{
				code:        202,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", nil)

			th := handlers.WrapperHandler{
				ServerConf: &config.ServiceShortURLConfig{
					Port:          8080,
					ServerHost:    "localhost",
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
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
			r.Delete("/api/user/urls", th.DeleteURLHandlers)

			r.ServeHTTP(w, request)

			//Code is checking result
			//I want a lot cod is checked
			//I will be  covering more tests by unit test in my code
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
