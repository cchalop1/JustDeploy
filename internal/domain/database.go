package domain

type DatabaseModelsType struct {
	Server   Server    `json:"server"`
	Services []Service `json:"services"`
	Settings Settings  `json:"settings"`
}
