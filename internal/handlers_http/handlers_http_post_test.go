package handlershttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/borisbbtest/go_home_work/internal/config"
	handlershttp "github.com/borisbbtest/go_home_work/internal/handlers_http"
	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func TestWrapperHandler_PostHandler(t *testing.T) {
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
				code:        201,
				response:    ``,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)

			th := handlershttp.WrapperHandler{
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
			h := http.HandlerFunc(th.PostHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// //resBody, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// // if string(resBody) != tt.want.response {
			// 	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			// }

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestWrapperHandler_PostJSONHandler(t *testing.T) {
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
			name: "Connetion test #2",
			want: want{
				code:        400,
				response:    "request body is not valid json",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", nil)

			th := handlershttp.WrapperHandler{
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
			r.Post("/api/shorten", th.PostJSONHandler)

			r.ServeHTTP(w, request)
			// запускаем сервер

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
			// resBody, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// if string(resBody) != tt.want.response {
			// 	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			// }

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestWrapperHandler_PostJSONHandlerBatch(t *testing.T) {
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
				code:        400,
				response:    `OK`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", nil)

			th := handlershttp.WrapperHandler{
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
			r.Post("/api/shorten/batch", th.PostJSONHandlerBatch)

			r.ServeHTTP(w, request)
			// запускаем сервер
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// // получаем и проверяем тело запроса
			defer res.Body.Close()
			// resBody, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// if string(resBody) != tt.want.response {
			// 	//	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			// }

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
