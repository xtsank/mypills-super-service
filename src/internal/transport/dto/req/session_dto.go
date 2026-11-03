package req

type SessionDto struct {
	Login    string `json:"login" required:"true"`
	Password string `json:"password" required:"true"`
}
