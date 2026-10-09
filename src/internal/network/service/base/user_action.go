package base

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type UserActionService struct{}

func (s *UserActionService) Create(req request.ActionCreateReq) error {
	userAction := basic.NewUserAction(req.UserID, req.TargetID, req.TargetType, req.ActionType)
	if err := global.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(userAction).Error; err != nil {
		global.Log.Error("用户操作失败", zap.Error(err))
		return fmt.Errorf("用户操作失败")
	}
	return nil
}
