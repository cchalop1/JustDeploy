package domain

type Server struct {
	Id          string
	Name        string
	Ip          string
	Domain      string
	Password    *string
	SshKey      *string
	CreatedDate string
}
