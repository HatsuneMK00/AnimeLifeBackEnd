package request

type Login struct {
	Type     LoginType `json:"type"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Code     string    `json:"code"`
}

type LoginType int

const (
	UsernamePassword LoginType = iota
	EmailVerificationCode
)
