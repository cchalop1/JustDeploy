package domain

import "cchalop1.com/deploy/internal/api/dto"

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
