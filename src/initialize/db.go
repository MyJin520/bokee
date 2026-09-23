package initialize

import (
	"bokee/config"
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/pkg/cryptox/hash"
	db "bokee/pkg/gormx"
	"bokee/pkg/randx"
	"fmt"
	"go.uber.org/zap"
)

func initDb(cfg *config.Config) {
	err := db.InitDB(cfg)
	if err != nil {
		panic("数据库初始化失败: " + err.Error())
	}
	migrateTable()
	initRolesAndUser()
}

// 迁移表结构
func migrateTable() {
	err := global.DB.AutoMigrate(
		&basic.User{},
		&basic.Role{},
		&basic.Article{},
		&basic.Files{},
		&basic.UserAction{},
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}

	// 历史遗留：手机号/邮箱为选填字段，旧表却带有 NOT NULL 唯一索引，
	// 多个未填手机号的用户会因空字符串冲突而无法注册；AutoMigrate 不会自动删除索引，这里显式移除。
	// 非空手机号/邮箱的唯一性由业务层校验保证。
	for _, indexName := range []string{"idx_users_phone", "idx_users_email"} {
		if global.DB.Migrator().HasIndex(&basic.User{}, indexName) {
			if err := global.DB.Migrator().DropIndex(&basic.User{}, indexName); err != nil {
				global.Log.Error(fmt.Sprintf("移除遗留唯一索引 %s 失败: %v", indexName, err))
			} else {
				global.Log.Info(fmt.Sprintf("已移除遗留唯一索引 %s", indexName))
			}
		}
	}
}

// 初始化内置角色与超级管理员：每次启动都会确保内置角色存在，并为历史无角色用户补授普通用户角色
func initRolesAndUser() {
	// 1. 确保内置角色存在（超级管理员、普通用户）
	adminRole := basic.Role{
		Name:   "超级管理员",
		Code:   uint(global.SuperRoleCode),
		Sort:   1,
		Status: "normal",
		Remark: "系统内置超级管理员角色",
	}
	if err := global.DB.Where("code = ?", adminRole.Code).FirstOrCreate(&adminRole).Error; err != nil {
		panic("初始化超级管理员角色失败: " + err.Error())
	}

	userRole := basic.Role{
		Name:   "普通用户",
		Code:   uint(global.UserRoleCode),
		Sort:   2,
		Status: "normal",
		Remark: "注册用户默认角色，拥有账户自助与个人文章管理权限",
	}
	if err := global.DB.Where("code = ?", userRole.Code).FirstOrCreate(&userRole).Error; err != nil {
		panic("初始化普通用户角色失败: " + err.Error())
	}

	// 2. 首次启动（用户表为空）时创建超级管理员
	var userCount int64
	if err := global.DB.Model(&basic.User{}).Count(&userCount).Error; err != nil {
		global.Log.Warn(fmt.Sprintf("检查用户表是否为空时出错: %v", err))
		return
	}
	if userCount == 0 {
		global.Log.Info("检测到空用户表，开始初始化超级管理员...")

		defaultPwd := global.Config.System.DefaultAdminPassword
		randomPwd := defaultPwd == ""
		if randomPwd {
			defaultPwd = randx.RandomDigitCode(10)
		}

		password, err := hash.GeneratePassword(defaultPwd)
		if err != nil {
			panic("密码加密失败: " + err.Error())
		}

		newUser := &basic.User{
			Name:     "superAdmin",
			Password: password,
			Status:   "normal",
			Roles:    []basic.Role{adminRole},
		}

		if err := global.DB.Create(newUser).Error; err != nil {
			panic("初始化超级管理员失败: " + err.Error())
		}

		global.Log.Info(fmt.Sprintf("✅ 初始化完成：已创建超级管理员账户（%s）并关联超级管理员角色", newUser.Name))
		if randomPwd {
			global.Log.Info(fmt.Sprintf("未配置 default-admin-password，已自动生成随机密码: %s，请登录后立即修改", defaultPwd))
		}
		return
	}

	global.Log.Info("用户表非空，跳过超级管理员初始化")
	backfillDefaultUserRole(userRole)
}

// backfillDefaultUserRole 为历史上没有任何角色的用户补授普通用户角色（注册逻辑上线前创建的账号）
func backfillDefaultUserRole(userRole basic.Role) {
	var users []basic.User
	if err := global.DB.Preload("Roles").Find(&users).Error; err != nil {
		global.Log.Warn(fmt.Sprintf("回填默认角色查询用户失败: %v", err))
		return
	}
	backfilled := 0
	for _, user := range users {
		if len(user.Roles) > 0 {
			continue
		}
		if err := global.DB.Model(&user).Association("Roles").Append(&userRole); err != nil {
			global.Log.Error("回填普通用户角色失败", zap.Uint("userID", user.ID), zap.Error(err))
			continue
		}
		backfilled++
	}
	if backfilled > 0 {
		global.Log.Info(fmt.Sprintf("已为 %d 个无角色用户补授普通用户角色", backfilled))
	}
}
