package domain

type ApplicationConfigResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ResponseApi struct {
	Message string `json:"message"`
}
