package models

type ProxyRequest struct {
	Method  string  `json:"method"`
	Url     string  `json:"url"`
	Headers *Header `json:"headers"`
}

type ProxyResponse struct {
	Id      string  `json:"id"`
	Status  string  `json:"status"`
	Headers *Header `json:"headers"`
	Length  int     `json:"int"`
}

type Header struct {
	Authentication string `json:"authentication"`
}
