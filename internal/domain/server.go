package domain

type Server struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Ip          string  `json:"ip"`
	Domain      string  `json:"domain"`
	Password    *string `json:"password"`
	SshKey      *string `json:"sshKey"`
	CreatedDate string  `json:"createdDate"`
	Status      string  `json:"status"`
}
