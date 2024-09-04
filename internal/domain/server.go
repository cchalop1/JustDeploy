package domain

import (
	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/dto"
)

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

type ServerCertsPath struct {
	CaCertPath string
	CertPath   string
	KeyPath    string
}

func (s *Server) ToServerDto() dto.ServerDto {
	return dto.ServerDto{
		Id:          s.Id,
		Name:        s.Name,
		Ip:          s.Ip,
		Domain:      s.Domain,
		CreatedDate: s.CreatedDate,
		Status:      s.Status,
	}
}

func (s *Server) GetCertsPath() ServerCertsPath {
	return ServerCertsPath{
		CaCertPath: internal.CERT_DOCKER_FOLDER + "/" + s.Id + "/ca.pem",
		CertPath:   internal.CERT_DOCKER_FOLDER + "/" + s.Id + "/cert.pem",
		KeyPath:    internal.CERT_DOCKER_FOLDER + "/" + s.Id + "/key.pem",
	}
}

func (s *Server) GetSshKeyPath() string {
	return internal.CERT_DOCKER_FOLDER + "/" + s.Id + "/"
}
