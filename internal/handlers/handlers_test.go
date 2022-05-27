package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/handlers"
	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/storage"
)

func TestGetShortLinkJSONHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		longURL     string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Create URL Short JSON test #2",
			want: want{
				code:        201,
				response:    `OK`,
				longURL:     "http://localhost:8080/.+",
				contentType: "application/json; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			v := model.RequestAddDBURL{
				ReqNewURL: "http://ya.ru",
			}
			reqBody, _ := json.Marshal(v)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", "application/json; charset=utf-8")
			th := handlers.WrapperHandler{
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
			h := http.HandlerFunc(th.PostJSONHandler)
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
			var m model.ResponseURLShort
			if err := json.Unmarshal(resBody, &m); err != nil {
				t.Errorf("Not valid JSON ")
			}
			if string(resBody) != tt.want.response {
				chk, _ := regexp.MatchString(tt.want.longURL, m.ResNewURL)
				if !chk {
					t.Errorf("Not valid url %s got %s ", tt.want.longURL, m.ResNewURL)
				}
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestGetShortLinkHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		longURL     string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Create URL Short test #1",
			want: want{
				code:        201,
				response:    `OK`,
				longURL:     "http://localhost:8080/.+",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("http://yandex.ru")))
			request.Header.Set("Content-Type", "text/plain; charset=utf-8")
			th := handlers.WrapperHandler{
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
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				chk, _ := regexp.MatchString(tt.want.longURL, w.Body.String())
				if !chk {
					t.Errorf("Not valid url ")
				}
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestStatusHandler(t *testing.T) {
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
				code:        200,
				response:    `OK`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)

			th := handlers.WrapperHandler{
				ServerConf: &config.ServiceShortURLConfig{
					Port:       8080,
					ServerHost: "localhost",
					BaseURL:    "http://localhost:8080",
				},
			}
			th.Storage, _ = &storage.NewFileStorage("")

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(th.GetHandler)
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

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
