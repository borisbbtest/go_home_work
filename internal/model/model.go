package model

type RequestAddDBURL struct {
	ReqNewURL string `json:"url"`
}
type ResponseURLShort struct {
	ResNewURL string `json:"result"`
}

type ResponseURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type ResponseStats struct {
	URLs  int32 `json:"urls"`
	Users int32 `json:"users"`
}
type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
type RequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type DataURL struct {
	Port          string `json:"Port"`
	URL           string `json:"URL"`
	Path          string `json:"Path"`
	ShortPath     string `json:"ShortPath"`
	UserID        string `json:"UserID"`
	CorrelationID string `json:"correlation_id"`
	StatusActive  int
}
