package base

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
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

func (s *UserActionService) Delete(actionId, userId uint) error {
	result := global.DB.Where("id = ? AND user_id = ?", actionId, userId).Delete(&basic.UserAction{})
	if result.Error != nil {
		global.Log.Error("用户操作删除失败", zap.Error(result.Error), zap.Uint("action_id", actionId), zap.Uint("user_id", userId))
		return fmt.Errorf("用户操作失败")
	}
	if result.RowsAffected == 0 {
		global.Log.Warn("用户操作删除失败：不存在或无权操作", zap.Uint("action_id", actionId), zap.Uint("user_id", userId))
		return fmt.Errorf("用户操作不存在")
	}
	return nil
}

func (s *UserActionService) List(req request.ActionListReq, userId uint) ([]response.UserActionListResponse, int64, error) {
	req.Normalize()

	switch req.TargetType {
	case "article":
		if req.ActionType != "like" && req.ActionType != "bookmark" {
			return nil, 0, fmt.Errorf("文章目标仅支持点赞或收藏")
		}
	case "author":
		if req.ActionType != "follow" {
			return nil, 0, fmt.Errorf("作者目标仅支持关注")
		}
	default:
		return nil, 0, fmt.Errorf("不支持的目标类型")
	}

	var total int64
	if err := global.DB.Model(&basic.UserAction{}).
		Where("user_id = ? AND target_type = ? AND action_type = ?", userId, req.TargetType, req.ActionType).
		Count(&total).Error; err != nil {
		global.Log.Error("统计用户操作记录数异常", zap.Error(err))
		return nil, 0, fmt.Errorf("统计用户操作记录数异常")
	}

	baseSelect := "ua.id AS action_id, ua.target_type, ua.target_id"
	query := global.DB.Table("user_actions AS ua").
		Where("ua.user_id = ? AND ua.target_type = ? AND ua.action_type = ?", userId, req.TargetType, req.ActionType).
		Order("ua.id DESC").
		Offset(req.Offset()).
		Limit(req.PageSize)

	switch req.TargetType {
	case "article":
		query = query.Select(baseSelect + ", t.title AS target_title").
			Joins("LEFT JOIN articles t ON t.id = ua.target_id")
	case "author":
		query = query.Select(baseSelect + ", t.name AS target_title").
			Joins("LEFT JOIN users t ON t.id = ua.target_id")
	}

	var list []response.UserActionListResponse
	if err := query.Scan(&list).Error; err != nil {
		global.Log.Error("查询用户操作记录列表异常", zap.Error(err))
		return nil, 0, fmt.Errorf("查询用户操作记录列表异常")
	}
	if list == nil {
		list = []response.UserActionListResponse{}
	}
	return list, total, nil
}
