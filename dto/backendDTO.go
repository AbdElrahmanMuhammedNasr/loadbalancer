package dto

type BackendDTO struct {
	NameSpace  string          `json:"name_space"`
	ServiceUrl []ServiceUrlDTO `json:"service_urls"`
}

type ServiceUrlDTO struct {
	Url    string `json:"url"`
	Active bool   `json:"active"`
}
