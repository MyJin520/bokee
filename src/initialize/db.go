package initialize

import (
	"errors"
	"fmt"
	"gin-admin/config"
	"gin-admin/global"
	"gin-admin/internal/mods/basic"
	"gin-admin/pkg/cryptox/hash"
	db "gin-admin/pkg/gormx"
	"gin-admin/pkg/randx"
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
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}

// 仅在首次启动时初始化角色和超级管理员
func initRolesAndUser() {
	var existingUser basic.User
	err := global.DB.Where("name = ?", "superAdmin").First(&existingUser).Error
	if err == nil {
		global.Log.Info("超级管理员已存在，跳过初始化")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Log.Warn(fmt.Sprintf("检查用户是否存在时出错: %v", err))
		return
	}

	global.Log.Info("检测到空数据库，开始初始化默认角色和管理员...")

	adminRole := basic.Role{
		RoleName: "超级管理员",
		RoleCode: 1000,
		Sort:     1,
		Status:   "normal",
		Remark:   "系统内置超级管理员角色",
	}
	result := global.DB.Where("role_code = ?", adminRole.RoleCode).FirstOrCreate(&adminRole)
	if result.Error != nil {
		panic("初始化默认角色失败: " + result.Error.Error())
	}

	defaultPwd := global.Config.System.DefaultAdminPassword
	if defaultPwd == "" {
		defaultPwd = randx.RandomDigitCode(10)
		global.Log.Info(fmt.Sprintf("默认超级管理员已创建： %s", existingUser.Name))
		global.Log.Info(fmt.Sprintf("未配置 default-admin-password，已自动生成随机密码: %s", defaultPwd))
		global.Log.Info(fmt.Sprintf("请登录后立即修改密码: %s", defaultPwd))
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

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newUser).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
		return nil
	})
	if err != nil {
		panic("初始化超级管理员失败: " + err.Error())
	}

	global.Log.Info("✅ 初始化完成：已创建超级管理员账户（superAdmin）并关联超级管理员角色")
}
