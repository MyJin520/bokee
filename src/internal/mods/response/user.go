package response

type JwtResp struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

type PriRouteResp struct {
	Path   string `json:"path"`   // 路由路径
	Method string `json:"method"` // 请求方法
}
