package tools

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("context", "service_short_url")

func AddCookie(w http.ResponseWriter, r *http.Request, name, value string, ttl time.Duration) (res string, err error) {
	ck, err := r.Cookie(name)
	if err != nil {
		log.Info("Cant find cookie : set cooke")
	} else {
		res = ck.Value
		return
	}
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
	return value, err
}

func GetCookie(r *http.Request, name string) (value string, err error) {
	cooke, err := r.Cookie(name)
	if err != nil {
		log.Error("Cant find cookie : set cooke")
		return
	}
	value = cooke.Value
	return
}
