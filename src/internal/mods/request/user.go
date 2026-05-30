package request

type UserRegisterOrEditReq struct {
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
