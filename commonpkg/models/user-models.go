package models

type UserDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	CommonResponse
	Token string `json:"token"`
}
