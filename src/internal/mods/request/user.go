package request

type UserRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenParsingReq struct {
	Token string `json:"token"`
}

type UserEditReq struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
