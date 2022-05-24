package model

type RequestAddDBURL struct {
	ReqNewURL string `json:"url"`
}
type ResponseURLShort struct {
	ResNewURL string `json:"result"`
}

type ResponseURLShortALL struct {
	ListsURL []ResponseURL `json:""`
}
type ResponseURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
