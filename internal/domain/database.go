package domain

type DatabaseModelsType struct {
	Servers  []Server  `json:"servers"`
	Deploys  []Deploy  `json:"deploys"`
	Services []Service `json:"services"`
	Projects []Project `json:"projects"` // Add this line
}
