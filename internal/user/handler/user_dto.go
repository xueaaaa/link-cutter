package handler

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
