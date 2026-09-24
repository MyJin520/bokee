package base

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
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

func (s *UserActionService) List(req request.ActionListReq, userId uint) ([]response.UserActionListResponse, int64, error) {
	req.Normalize()
	query := global.DB.Model(&basic.UserAction{}).Where("user_id = ?", userId).Where("action_type = ?", req.ActionType)
	if req.TargetType != "" {
		query = query.Where("target_type = ?", req.TargetType)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Log.Error("用户操作列表获取失败", zap.Error(err))
		return nil, 0, fmt.Errorf("用户操作列表获取失败")
	}

	var actions []basic.UserAction
	if err := query.Order("id DESC").Offset(req.Offset()).Limit(req.PageSize).Find(&actions).Error; err != nil {
		global.Log.Error("用户操作列表查询失败", zap.Error(err))
		return nil, 0, fmt.Errorf("用户操作列表获取失败")
	}

	list := make([]response.UserActionListResponse, 0, len(actions))
	for _, action := range actions {
		item := response.UserActionListResponse{ActionID: action.ID}
		switch action.TargetType {
		case "article":
			var article basic.Article
			if err := global.DB.Select("id", "title").Where("id = ?", action.TargetID).First(&article).Error; err != nil {
				global.Log.Warn("查询目标文章失败", zap.Error(err), zap.Uint("article_id", action.TargetID))
				continue
			}
			item.Article = &response.Article{ID: article.ID, Title: article.Title}
		case "author":
			var user basic.User
			if err := global.DB.Select("id", "name").Where("id = ?", action.TargetID).First(&user).Error; err != nil {
				global.Log.Warn("查询目标作者失败", zap.Error(err), zap.Uint("author_id", action.TargetID))
				continue
			}
			item.Author = &response.Author{ID: user.ID, Name: user.Name}
		}
		list = append(list, item)
	}
	return list, total, nil
}
