package handlershttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/borisbbtest/go_home_work/internal/tools"
)

// AddCookie добавим мидлу
func AddCookie(w http.ResponseWriter, r *http.Request, name, value string, ttl time.Duration) (res string, err error) {
	ck, err := r.Cookie(name)
	if err == nil {
		res = ck.Value
		return
	}
	log.Info("Cant find cookie : set cooke")
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
	return value, err
}

func (hook *WrapperHandler) MidSetCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmp, _ := tools.GetID()
		hook.UserID, _ = AddCookie(w, r, "ShortURL", fmt.Sprintf("%x", tmp), 30*time.Minute)
		next.ServeHTTP(w, r)
	})
}
