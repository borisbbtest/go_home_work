package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
)

func (hook *WrapperHandler) DeleteURLHandlers(w http.ResponseWriter, r *http.Request) {

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}

	log.Info("PostHandler -> ", string(bytes))

	log.Println("Post handler")
	defer r.Body.Close()
}
