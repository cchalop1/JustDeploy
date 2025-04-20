package domain

import (
	"time"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/api/dto"
)

type Server struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	Port        string    `json:"port"`
	Domain      string    `json:"domain"`
	Password    *string   `json:"password"`
	SshKey      *string   `json:"sshKey"`
	CreatedDate time.Time `json:"createdDate"`
	Status      string    `json:"status"`
	Email       string    `json:"email"`
	UseHttps    bool      `json:"useHttps"`
}

type ServerCertsPath struct {
	CaCertPath string
	CertPath   string
	KeyPath    string
}

func (s *Server) ToServerDto() dto.ServerDto {
	return dto.ServerDto{
		Id:       s.Id,
		Name:     s.Name,
		Ip:       s.Ip,
		Domain:   s.Domain,
		Status:   s.Status,
		UseHttps: s.UseHttps,
		Email:    s.Email,
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
