package initialize

import (
	"bokee/config"
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/pkg/cryptox/hash"
	db "bokee/pkg/gormx"
	"bokee/pkg/randx"
	"fmt"
	"gorm.io/gorm"
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
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}

// 仅在首次启动时初始化角色和超级管理员
func initRolesAndUser() {
	var userCount int64
	if err := global.DB.Model(&basic.User{}).Count(&userCount).Error; err != nil {
		global.Log.Warn(fmt.Sprintf("检查用户表是否为空时出错: %v", err))
		return
	}
	if userCount > 0 {
		global.Log.Info("用户表非空，跳过超级管理员初始化")
		return
	}

	global.Log.Info("检测到空用户表，开始初始化默认角色和管理员...")

	adminRole := basic.Role{
		RoleName: "超级管理员",
		RoleCode: uint(global.SuperRoleCode),
		Sort:     1,
		Status:   "normal",
		Remark:   "系统内置超级管理员角色",
	}

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
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_code = ?", adminRole.RoleCode).FirstOrCreate(&adminRole).Error; err != nil {
			return fmt.Errorf("初始化默认角色失败: %w", err)
		}
		newUser.Roles = []basic.Role{adminRole}
		if err := tx.Create(newUser).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
		return nil
	})
	if err != nil {
		panic("初始化超级管理员失败: " + err.Error())
	}

	global.Log.Info(fmt.Sprintf("✅ 初始化完成：已创建超级管理员账户（%s）并关联超级管理员角色", newUser.Name))
	if randomPwd {
		global.Log.Info(fmt.Sprintf("未配置 default-admin-password，已自动生成随机密码: %s，请登录后立即修改", defaultPwd))
	}
}
