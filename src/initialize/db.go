package initialize

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gin-admin/config"
	"gin-admin/global"
	"gin-admin/internal/mods/basic"
	"gin-admin/pkg/cryptox/hash"
	db "gin-admin/pkg/gormx"
)

func initDb(cfg *config.Config) {
	err := db.InitDB(cfg)
	if err != nil {
		panic("数据库初始化失败: " + err.Error())
	}
	migrateTable()
	initUser()
}

// 迁移表结构
func migrateTable() {
	err := global.DB.AutoMigrate(
		// 基础模型
		&basic.User{},
		&basic.Role{},
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}

func initUser() {
	var count int64
	global.DB.Model(&basic.User{}).Count(&count)
	if count > 0 {
		return
	}

	// 优先使用配置文件中的密码，留空则自动生成随机密码
	defaultPwd := global.Config.System.DefaultAdminPassword
	if defaultPwd == "" {
		defaultPwd = generateRandomPassword()
		fmt.Println()
		fmt.Println("╔══════════════════════════════════════════════════════════╗")
		fmt.Println("║  未配置 default-admin-password，已自动生成随机密码        ║")
		fmt.Println("║                                                          ║")
		fmt.Printf("║  用户名:  superAdmin                                     ║\n")
		fmt.Printf("║  密  码:  %-46s ║\n", defaultPwd)
		fmt.Println("║                                                          ║")
		fmt.Println("║  请登录后立即修改！                                       ║")
		fmt.Println("╚══════════════════════════════════════════════════════════╝")
		fmt.Println()
	}

	password, err := hash.GeneratePassword(defaultPwd)
	if err != nil {
		panic("密码加密失败: " + err.Error())
	}

	user := &basic.User{
		Name:     "superAdmin",
		Password: password,
		Status:   "normal",
	}

	// FirstOrCreate 避免多实例同时启动时唯一约束冲突导致 panic
	result := global.DB.Where(basic.User{Name: "superAdmin"}).FirstOrCreate(user)
	if result.Error != nil {
		panic("初始化默认管理员失败: " + result.Error.Error())
	}
	if result.RowsAffected > 0 {
		global.Log.Info("已创建默认管理员账户（superAdmin）")
	}
}

// generateRandomPassword 生成 22 位安全随机密码
func generateRandomPassword() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
