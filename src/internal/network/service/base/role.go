package base

import (
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/pkg/casbinx"
)

type RoleService struct{}

func (s *RoleService) GetAllPriRoles(page request.PageReq) ([]response.PriRouteResp, int64, error) {
	// 从 Casbin 缓存/表读取私有路由策略
	policies, err := casbinx.GetPrivateRoutes()
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应结构体
	var routes []response.PriRouteResp
	for _, p := range policies {
		routes = append(routes, response.PriRouteResp{
			Path:   p[1], // obj
			Method: p[2], // act
		})
	}

	// 内存分页
	total := int64(len(routes))
	start := page.Offset()
	if start > int(total) {
		return []response.PriRouteResp{}, total, nil
	}
	end := start + page.PageSize
	if end > int(total) {
		end = int(total)
	}

	return routes[start:end], total, nil
}
