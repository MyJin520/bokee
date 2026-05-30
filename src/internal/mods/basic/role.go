package basic

type Role struct {
	BaseModel
	RoleName string `gorm:"size:20;not null;uniqueIndex;comment:角色名称"`
	RoleCode string `gorm:"size:20;not null;uniqueIndex;comment:角色标识"`
	Sort     int    `gorm:"not null;default:0;comment:排序序号"`
	Status   string `gorm:"size:10;not null;default:normal;comment:角色状态"`
	Remark   string `gorm:"size:255;comment:角色备注说明"`

	Users []User `gorm:"many2many:sys_user_roles;constraint:OnDelete:CASCADE"`
}
