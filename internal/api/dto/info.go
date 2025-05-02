package dto

type Info struct {
	Version         VersionDto `json:"version"`
	FirstConnection bool       `json:"firstConnection"`
	Server          ServerDto  `json:"server"`
}
