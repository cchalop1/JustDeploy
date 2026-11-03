package dto

type InitialSetupDto struct {
	SetupToken string `json:"setupToken"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Domain     string `json:"domain"`
}
