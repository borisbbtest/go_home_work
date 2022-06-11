package handlers

import (
	"compress/gzip"
	"encoding/json"
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

	bytesBody, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}

	var shortURLs []string
	json.Unmarshal(bytesBody, &shortURLs)

	log.Info("DeletedHandler -> ")
	log.Info(shortURLs)

	w.WriteHeader(http.StatusAccepted)
	log.Println("Deleted handler")
	defer r.Body.Close()
}
