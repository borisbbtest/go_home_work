package model

type RequestAddDBURL struct {
	ReqNewURL string `json:"url"`
}
type ResponseURLShort struct {
	ResNewURL string `json:"result"`
}
