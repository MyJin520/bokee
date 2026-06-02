package casbinx

import (
	"fmt"
	"gin-admin/global"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"sync"
)

var (
	syncedEnforcer *casbin.SyncedCachedEnforcer
	once           sync.Once
	initErr        error
)

// casbinModelText 定义 RBAC 模型：
// - g(r.sub, p.sub) 支持角色继承（用户→角色→权限）
// - keyMatch3 支持 RESTful 路径匹配（/api/user/:id）
// - regexMatch 支持正则匹配 HTTP 方法
const casbinModelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch3(r.obj, p.obj) && regexMatch(r.act, p.act)
`

// InitEnforcer 初始化 Casbin 执行器（程序启动时调用一次）。
// 失败时返回 error，调用方应阻止程序继续启动。
func InitEnforcer() error {
	once.Do(func() {
		// 1. 创建 GORM 适配器（策略存储在 casbin_rule 表中）
		adapter, err := gormadapter.NewAdapterByDB(global.DB)
		if err != nil {
			initErr = fmt.Errorf("casbin 适配器创建失败: %w", err)
			global.Log.Error("casbin 适配器创建失败", zap.Error(err))
			return
		}

		// 2. 从字符串解析 RBAC 模型
		m, err := model.NewModelFromString(casbinModelText)
		if err != nil {
			initErr = fmt.Errorf("casbin 模型解析失败: %w", err)
			global.Log.Error("casbin 模型解析失败", zap.Error(err))
			return
		}

		// 3. 创建带缓存和并发安全的执行器
		enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
		if err != nil {
			initErr = fmt.Errorf("casbin 执行器创建失败: %w", err)
			global.Log.Error("casbin 执行器创建失败", zap.Error(err))
			return
		}

		// 4. 设置缓存过期时间（秒），过期后自动从数据库重新加载策略
		enforcer.SetExpireTime(60 * 60)

		// 5. 启动时预加载所有策略到内存缓存
		if err := enforcer.LoadPolicy(); err != nil {
			initErr = fmt.Errorf("casbin 策略加载失败: %w", err)
			global.Log.Error("casbin 策略加载失败", zap.Error(err))
			return
		}

		syncedEnforcer = enforcer
		global.Log.Info("casbin 初始化完成")
	})
	return initErr
}

// GetEnforcer 返回已初始化的 Casbin 执行器。
// 注意：必须先调用 InitEnforcer 并确保成功，否则返回 nil 和对应的错误。
func GetEnforcer() (*casbin.SyncedCachedEnforcer, error) {
	if syncedEnforcer == nil {
		return nil, initErr
	}
	return syncedEnforcer, nil
}

// Enforce 权限校验。
// sub: 用户名，obj: 请求路径（如 /pri/user/info），act: HTTP 方法（GET/POST/PUT/DELETE）。
func Enforce(sub, obj, act string) (bool, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return false, err
	}
	return enforcer.Enforce(sub, obj, act)
}

// AddPolicyForUser 为用户直接添加权限策略。
func AddPolicyForUser(user, path, method string) (bool, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return false, err
	}
	return enforcer.AddPolicy(user, path, method)
}

// RemovePolicyForUser 移除用户的直接权限策略。
func RemovePolicyForUser(user, path, method string) (bool, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return false, err
	}
	return enforcer.RemovePolicy(user, path, method)
}

// AddRoleForUser 为用户分配角色（写入 g 规则）。
func AddRoleForUser(user, role string) (bool, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return false, err
	}
	return enforcer.AddGroupingPolicy(user, role)
}

// RemoveRoleForUser 移除用户的角色。
func RemoveRoleForUser(user, role string) (bool, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return false, err
	}
	return enforcer.RemoveGroupingPolicy(user, role)
}

// GetRolesForUser 获取用户的直接角色（不包含继承）。
func GetRolesForUser(user string) ([]string, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return nil, err
	}
	return enforcer.GetRolesForUser(user)
}

// GetImplicitRolesForUser 获取用户的所有角色（包含继承链路）。
func GetImplicitRolesForUser(user string) ([]string, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return nil, err
	}
	return enforcer.GetImplicitRolesForUser(user)
}

// ReloadPolicy 从数据库重新加载所有策略到内存缓存。
// 管理员通过 API 修改权限后，调用此方法使其立即生效。
func ReloadPolicy() error {
	enforcer, err := GetEnforcer()
	if err != nil {
		return err
	}
	return enforcer.LoadPolicy()
}

// GetAllPolicies 获取所有权限策略（p 规则）。
func GetAllPolicies() ([][]string, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return nil, err
	}
	return enforcer.GetPolicy()
}

// GetAllRoles 获取所有角色分配关系（g 规则）。
func GetAllRoles() ([][]string, error) {
	enforcer, err := GetEnforcer()
	if err != nil {
		return nil, err
	}
	return enforcer.GetGroupingPolicy()
}
