package handlershttp

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"github.com/borisbbtest/go_home_work/internal/model"
)

// Хедлер удаления
func (hook *WrapperHandler) DeleteURLHandlers(w http.ResponseWriter, r *http.Request) {

	// nr
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
	//пока формально Нет времени подумать едином канале
	go func() {
		buff := make([]model.DataURL, 0, len(shortURLs))
		for _, str := range shortURLs {
			buff = append(buff, model.DataURL{ShortPath: str,
				StatusActive: 2})
		}
		hook.Storage.DeletedURLBatch(hook.UserID, buff)
	}()

	log.Info("DeletedHandler -> ")
	log.Info(shortURLs)

	w.WriteHeader(http.StatusAccepted)
	log.Println("Deleted handler")
	defer r.Body.Close()
}
