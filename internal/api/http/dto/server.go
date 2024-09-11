package dto

type ConnectNewServerDto struct {
	Domain   string  `json:"domain"`
	Ip       string  `json:"ip"`
	SshKey   *string `json:"sshKey"`
	Password *string `json:"password"`
	User     string  `json:"user"`
}

type ServerDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	Domain      string `json:"domain"`
	CreatedDate string `json:"createdDate"`
	Status      string `json:"status"`
}
