package dto

type ConnectNewServerDto struct {
	Domain   string  `json:"domain"`
	Ip       string  `json:"ip"`
	Email    string  `json:"email"`
	SshKey   *string `json:"sshKey"`
	Password *string `json:"password"`
	User     string  `json:"user"`
}

type ServerDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Ip       string `json:"ip"`
	Domain   string `json:"domain"`
	Status   string `json:"status"`
	UseHttps bool   `json:"useHttps"`
	Email    string `json:"email"`
}
