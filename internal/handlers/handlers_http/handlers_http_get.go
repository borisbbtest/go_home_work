package handlershttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/borisbbtest/go_home_work/internal/storage"
	"github.com/borisbbtest/go_home_work/internal/tools"
	"github.com/go-chi/chi/v5"
)

// GetHandler Запрос короткого URL
func (hook *WrapperHandler) GetHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	// for k, v := range hook.UrlStore.DBLocal {
	// 	fmt.Printf("key[%s] value[%s]\n", k, v.Url)
	// }
	log.Info("ID Go to", id)
	value, status := hook.Storage.Get(id)
	if status == nil {
		if value.StatusActive == storage.RecordNotFound {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusGone)
			fmt.Fprint(w, "Short url deleted")
			return
		}
		url := value.URL
		w.Header().Set("Location", url)
		w.WriteHeader(307)
		//log.Printf("Get handler")
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		fmt.Fprint(w, "OK")
		//log.Printf("key not found")
	}
	defer r.Body.Close()
	//log.Printf("Get handler")
}

// GetHandlerCooke Настройка кукизов
func (hook *WrapperHandler) GetHandlerCooke(w http.ResponseWriter, r *http.Request) {

	if len(hook.UserID) == 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(204)
		fmt.Fprint(w, "No Content")
		return
		//log.Printf("Get handler")
	}

	responseShortURL, err := hook.Storage.GetAll(hook.UserID, hook.ServerConf.BaseURL)

	if err != nil || len(responseShortURL) <= 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(204)
		fmt.Fprint(w, "No Content")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(responseShortURL)

	//log.Printf("Get handler")
}

func (hook *WrapperHandler) GetHandlerStats(w http.ResponseWriter, r *http.Request) {

	if len(hook.UserID) == 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "No Content")
		return
		//log.Printf("Get handler")
	}

	_, err := tools.TrustedSubnet(r.Header.Get(r.Header.Get("X-Real-IP")), *hook.ServerConf.Subnet)
	if err != nil {
		ips := r.Header.Get("X-Forwarded-For")
		ipStrs := strings.Split(ips, ",")
		ipStr := ipStrs[0]
		_, err = tools.TrustedSubnet(ipStr, *hook.ServerConf.Subnet)
	}
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		log.Error(err)
		fmt.Fprint(w, "403 Forbidden")
		return
	}

	responseShortURL, _ := hook.Storage.GetStats()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(responseShortURL)

	//log.Printf("Get handler")
}

// GetHandlerPing пинголятор
func (hook *WrapperHandler) GetHandlerPing(w http.ResponseWriter, r *http.Request) {

	_, status := tools.PingDataBase(hook.ServerConf.DataBaseDSN)
	if status != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(500)
		fmt.Fprint(w, "error connection")
		return
		//log.Printf("Get handler")
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprint(w, "ok")

	//log.Printf("Get handler")
}
