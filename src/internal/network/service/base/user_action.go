package base

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type UserActionService struct{}

func (s *UserActionService) Create(req request.ActionCreateReq) error {
	userAction := basic.NewUserAction(req.UserID, req.TargetID, req.TargetType, req.ActionType)
	if err := global.DB.Create(userAction).Error; err != nil {
		global.Log.Error("用户操作创建失败", zap.Error(err))
		return fmt.Errorf("用户操作创建失败")
	}
	return nil
}

func (s *UserActionService) Delete(ctx context.Context, actionId uint, userId uint) error {
	result := global.DB.Where("id = ? AND user_id = ?", actionId, userId).Delete(&basic.UserAction{})
	if result.Error != nil {
		global.Log.Error("用户操作删除失败", zap.Error(result.Error))
		return fmt.Errorf("用户操作删除失败")
	}
	if result.RowsAffected == 0 {
		global.Log.Warn("用户操作删除失败：用户操作不存在或无权操作", zap.Uint("user_id", userId), zap.Uint("action_id", actionId))
		return fmt.Errorf("用户操作不存在")
	}
	global.Log.Info("用户操作删除成功", zap.Uint("user_id", userId), zap.Uint("action_id", actionId))
	return nil
}
