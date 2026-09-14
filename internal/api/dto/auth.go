package dto

type RegisterDto struct {
	Token    string `json:"token" validate:"required"`
	Login    string `json:"login" validate:"required,min=8,alphanum"`
	Password string `json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789,containsany=!@#$%^&*()_+-=~"`
}

type AuthDto struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
