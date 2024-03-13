package dto

type ConnectServerDto struct {
	Domain   string  `json:"domain"`
	SshKey   *string `json:"sshKey"`
	Password *string `json:"password"`
	User     string  `json:"user"`
}
