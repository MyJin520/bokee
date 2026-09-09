package base

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/pkg/casbinx"
	"errors"
	"go.ube
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type RoleService struct{}

// Create 创建角色
func (s *RoleService) Create(req request.RoleCreateReq) (*response.RoleResp, error) {
	// 校验角色名称是否已存在
	var nameCount int64
	global.DB.Model(&basic.Role{}).Where("role_name = ?", req.RoleName).Count(&nameCount)
	if nameCount > 0 {
		return nil, errors.New("角色名称已存在")
	}

	// 校验角色标识是否已存在
	var codeCount int64
	global.DB.Model(&basic.Role{}).Where("role_code = ?", req.RoleCode).Count(&codeCount)
	if codeCount > 0 {
		return nil, errors.New("角色标识已存在")
	}

	// 默认值
	status := req.Status
	if status == "" {
		status = "normal"
	}

	role := &basic.Role{
		RoleName: req.RoleName,
		RoleCode: req.RoleCode,
		Sort:     req.Sort,
		Status:   status,
		Remark:   req.Remark,
	}

	if err := global.DB.Create(role).Error; err != nil {
		global.Log.Error("创建角色失败", zap.Error(err))
		return nil, errors.New("创建角色失败，请稍后重试")
	}

	return &response.RoleResp{
		ID:        role.ID,
		RoleName:  role.RoleName,
		RoleCode:  role.RoleCode,
		Sort:      role.Sort,
		Status:    role.Status,
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}, nil
}

// Update 更新角色
func (s *RoleService) Update(req request.RoleUpdateReq) error {
	// 查询原角色
	var role basic.Role
	err := global.DB.Where("id = ?", req.ID).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		global.Log.Error("查询角色失败", zap.Error(err))
		return errors.New("查询角色失败，请稍后重试")
	}

	// 禁止修改超级管理员角色信息
	if role.RoleCode == uint(global.SuperRoleCode) {
		return errors.New("超级管理员角色不可修改")
	}

	// 校验角色名称是否与其他角色冲突
	if req.RoleName != "" && req.RoleName != role.RoleName {
		var nameCount int64
		global.DB.Model(&basic.Role{}).Where("role_name = ? AND id != ?", req.RoleName, req.ID).Count(&nameCount)
		if nameCount > 0 {
			return errors.New("角色名称已存在")
		}
	}
	// todo 后续使用公共更新方法
	// 构建更新字段
	updates := make(map[string]interface{})
	if req.RoleName != "" {
		updates["role_name"] = req.RoleName
	}
	if req.Sort != role.Sort {
		updates["sort"] = req.Sort
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	updates["remark"] = req.Remark // remark 允许置空

	if len(updates) == 0 {
		return errors.New("没有需要更新的字段")
	}

	result := global.DB.Model(&basic.Role{}).Where("id = ?", req.ID).Updates(updates)
	if result.Error != nil {
		global.Log.Error("更新角色失败", zap.Error(result.Error))
		return errors.New("更新角色失败，请稍后重试")
	}
	if result.RowsAffected == 0 {
		return errors.New("角色不存在或未做任何更改")
	}

	return nil
}

// Delete 删除角色（软删除）
func (s *RoleService) Delete(id uint) error {
	// 查询角色
	var role basic.Role
	err := global.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		global.Log.Error("查询角色失败", zap.Error(err))
		return errors.New("查询角色失败，请稍后重试")
	}

	// 禁止删除超级管理员角色
	if role.RoleCode == uint(global.SuperRoleCode) {
		return errors.New("超级管理员角色不可删除")
	}

	// 检查是否有用户关联此角色
	var userCount int64
	global.DB.Table("sys_user_roles").Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		return errors.New("该角色下存在用户，无法删除，请先解除用户关联")
	}

	if err := global.DB.Delete(&role).Error; err != nil {
		global.Log.Error("删除角色失败", zap.Error(err), zap.Uint("roleID", id))
		return errors.New("删除角色失败，请稍后重试")
	}

	global.Log.Info("删除角色成功", zap.Uint("roleID", id), zap.String("roleName", role.RoleName))
	return nil
}

// GetInfo 获取单个角色详情
func (s *RoleService) GetInfo(id uint) (*response.RoleResp, error) {
	var role basic.Role
	err := global.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		global.Log.Error("查询角色失败", zap.Error(err))
		return nil, errors.New("查询角色失败，请稍后重试")
	}
	// todo 后续直接使用表结构体本身
	return &response.RoleResp{
		ID:        role.ID,
		RoleName:  role.RoleName,
		RoleCode:  role.RoleCode,
		Sort:      role.Sort,
		Status:    role.Status,
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}, nil
}

// List 分页获取角色列表
func (s *RoleService) List(req request.RoleQueryReq) ([]response.RoleResp, int64, error) {
	// 构建查询
	query := global.DB.Model(&basic.Role{})

	// 可选过滤条件
	if req.RoleName != "" {
		query = query.Where("role_name LIKE ?", "%"+req.RoleName+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Log.Error("统计角色总数失败", zap.Error(err))
		return nil, 0, errors.New("查询角色列表失败，请稍后重试")
	}

	// 分页查询
	var roles []basic.Role
	if err := query.Order("sort ASC, id ASC").
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&roles).Error; err != nil {
		global.Log.Error("查询角色列表失败", zap.Error(err))
		return nil, 0, errors.New("查询角色列表失败，请稍后重试")
	}
	// todo 后续优化
	// 转换为响应结构体
	list := make([]response.RoleResp, 0, len(roles))
	for _, role := range roles {
		list = append(list, response.RoleResp{
			ID:        role.ID,
			RoleName:  role.RoleName,
			RoleCode:  role.RoleCode,
			Sort:      role.Sort,
			Status:    role.Status,
			Remark:    role.Remark,
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
		})
	}

	return list, total, nil
}

// Auth 为角色批量授权（添加 Casbin 策略）
func (s *RoleService) Auth(req request.RoleAuthReq) error {
	// 查询角色
	var role basic.Role
	err := global.DB.Where("id = ?", req.RoleID).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		global.Log.Error("查询角色失败", zap.Error(err))
		return errors.New("查询角色失败，请稍后重试")
	}

	// 禁止为超级管理员手动授权（超级管理员自动拥有全部权限）
	if role.RoleCode == uint(global.SuperRoleCode) {
		return errors.New("超级管理员无需手动授权")
	}

	if len(req.Rules) == 0 {
		return errors.New("授权规则不能为空")
	}

	// 构建 Casbin 策略：p = sub(roleCode), obj(path), act(method)
	roleCodeStr := strconv.Itoa(int(role.RoleCode))
	policies := make([][]string, 0, len(req.Rules))
	for _, rule := range req.Rules {
		if rule.Path == "" || rule.Method == "" {
			return errors.New("路由路径和请求方法不能为空")
		}
		policies = append(policies, []string{roleCodeStr, rule.Path, strings.ToUpper(rule.Method)})
	}

	enforcer, err := casbinx.GetEnforcer()
	if err != nil {
		global.Log.Error("获取 Casbin 执行器失败", zap.Error(err))
		return errors.New("授权失败，请稍后重试")
	}

	// 批量添加策略
	success, err := enforcer.AddPolicies(policies)
	if err != nil {
		global.Log.Error("批量添加 Casbin 策略失败", zap.Error(err))
		return errors.New("授权失败，请稍后重试")
	}
	if !success {
		return errors.New("授权规则与已有规则重复，无需重复添加")
	}

	// 重载策略使更改立即生效
	if err := casbinx.ReloadPolicy(); err != nil {
		global.Log.Error("重载 Casbin 策略失败", zap.Error(err))
		return errors.New("授权成功但策略重载失败，请稍后重试")
	}

	global.Log.Info("角色授权成功",
		zap.Uint("roleID", req.RoleID),
		zap.String("roleName", role.RoleName),
		zap.Int("ruleCount", len(req.Rules)),
	)
	return nil
}

// GetAllPriRule 获取所有私有路由（原有功能，保留兼容）
func (s *RoleService) GetAllPriRule(page request.PageReq) ([]response.PriRouteResp, int64, error) {
	policies, err := casbinx.GetPrivateRoutes()
	if err != nil {
		return nil, 0, err
	}

	var routes []response.PriRouteResp
	for _, p := range policies {
		routes = append(routes, response.PriRouteResp{
			Path:   p[1],
			Method: p[2],
		})
	}

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
