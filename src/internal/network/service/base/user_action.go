package base

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/pkg/cachex"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserActionService struct{}

// Create 创建用户操作（点赞/收藏/关注）
// 幂等处理，避免唯一索引 (user_id, target_id, target_type, action_type) 冲突：
// 1. 已有有效记录 → 直接视为成功（重复操作不报错）；
// 2. 存在软删除残留行 → 复活原记录（取消后再次操作，恢复旧行而非新建）；
// 3. 均不存在 → 新建（并发兜底：唯一索引冲突时静默成功）。
// 点赞文章时同步自增文章点赞数（仅在实际新建/复活时）。
func (s *UserActionService) Create(req request.ActionCreateReq) error {
	created := false
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 幂等：有效记录已存在时视为成功
		var activeCount int64
		if err := tx.Model(&basic.UserAction{}).
			Where("user_id = ? AND target_id = ? AND target_type = ? AND action_type = ?",
				req.UserID, req.TargetID, req.TargetType, req.ActionType).
			Count(&activeCount).Error; err != nil {
			return err
		}
		if activeCount > 0 {
			return nil
		}

		// 2. 软删除残留行复活（Unscoped 才能命中已删除记录）
		resurrect := tx.Model(&basic.UserAction{}).Unscoped().
			Where("user_id = ? AND target_id = ? AND target_type = ? AND action_type = ? AND deleted_at IS NOT NULL",
				req.UserID, req.TargetID, req.TargetType, req.ActionType).
			Update("deleted_at", nil)
		if resurrect.Error != nil {
			return resurrect.Error
		}
		if resurrect.RowsAffected > 0 {
			created = true
			return nil
		}

		// 3. 无残留行时新建；并发下唯一索引冲突时静默成功（不重复创建）
		userAction := basic.NewUserAction(req.UserID, req.TargetID, req.TargetType, req.ActionType)
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(userAction).Error; err != nil {
			return err
		}
		if userAction.ID > 0 {
			created = true
		}
		return nil
	})
	if err != nil {
		global.Log.Error("用户操作创建失败", zap.Error(err))
		return fmt.Errorf("用户操作创建失败")
	}

	// 点赞文章时同步自增文章点赞数（取消点赞时已扣减，再次点赞补回）
	if created && req.ActionType == "like" && req.TargetType == "article" {
		syncArticleLikeCount(req.TargetID, 1)
	}
	return nil
}

// Delete 删除用户操作；取消点赞文章时同步扣减文章点赞数
func (s *UserActionService) Delete(ctx context.Context, actionId uint, userId uint) error {
	// 删除前取回操作记录，用于取消点赞时同步扣减文章点赞数
	var action basic.UserAction
	if err := global.DB.Where("id = ? AND user_id = ?", actionId, userId).First(&action).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("用户操作删除失败：用户操作不存在或无权操作", zap.Uint("user_id", userId), zap.Uint("action_id", actionId))
			return fmt.Errorf("用户操作不存在")
		}
		global.Log.Error("用户操作查询失败", zap.Error(err), zap.Uint("user_id", userId), zap.Uint("action_id", actionId))
		return fmt.Errorf("用户操作删除失败")
	}

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

	// 取消点赞文章时同步扣减文章点赞数
	if action.ActionType == "like" && action.TargetType == "article" {
		syncArticleLikeCount(action.TargetID, -1)
	}
	return nil
}

// syncArticleLikeCount 点赞/取消点赞文章时同步文章点赞数字段（GREATEST 防止负值）。
// 同步失败仅记录日志，不阻断主流程；成功后失效文章详情缓存避免展示过期。
func syncArticleLikeCount(articleID uint, delta int) {
	if err := global.DB.Model(&basic.Article{}).Where("id = ?", articleID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count + ?, 0)", delta)).Error; err != nil {
		global.Log.Warn("文章点赞数同步失败", zap.Uint("articleID", articleID), zap.Int("delta", delta), zap.Error(err))
		return
	}
	cachex.Invalidate(context.Background(), cachex.Key(cachex.NSArticle, "info", articleID))
}

func (s *UserActionService) List(req request.ActionListReq, userID uint) ([]response.UserActionListResponse, int64, error) {
	req.Normalize()

	query := global.DB.Model(&basic.UserAction{}).Where("user_id = ?", userID).Where("action_type = ?", req.ActionType)
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

	articleIDs := collectTargetIDs(actions, "article")
	authorIDs := collectTargetIDs(actions, "author")

	articleMap := make(map[uint]*response.Article, len(articleIDs))
	if len(articleIDs) > 0 {
		var articles []basic.Article
		if err := global.DB.Select("id", "title").Where("id IN ?", articleIDs).Find(&articles).Error; err != nil {
			global.Log.Error("批量查询文章失败", zap.Error(err))
			return nil, 0, fmt.Errorf("用户操作列表获取失败")
		}
		for _, article := range articles {
			articleMap[article.ID] = &response.Article{ID: article.ID, Title: article.Title}
		}
	}

	authorMap := make(map[uint]*response.Author, len(authorIDs))
	if len(authorIDs) > 0 {
		var users []basic.User
		if err := global.DB.Select("id", "name").Where("id IN ?", authorIDs).Find(&users).Error; err != nil {
			global.Log.Error("批量查询作者失败", zap.Error(err))
			return nil, 0, fmt.Errorf("用户操作列表获取失败")
		}
		for _, user := range users {
			authorMap[user.ID] = &response.Author{ID: user.ID, Name: user.Name}
		}
	}

	list := make([]response.UserActionListResponse, 0, len(actions))
	for _, action := range actions {
		item := response.UserActionListResponse{ActionID: action.ID}
		switch action.TargetType {
		case "article":
			item.Article = articleMap[action.TargetID]
		case "author":
			item.Author = authorMap[action.TargetID]
		}
		list = append(list, item)
	}

	return list, total, nil
}

func collectTargetIDs(actions []basic.UserAction, targetType string) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0)
	for _, action := range actions {
		if action.TargetType != targetType {
			continue
		}
		if _, ok := seen[action.TargetID]; ok {
			continue
		}
		seen[action.TargetID] = struct{}{}
		ids = append(ids, action.TargetID)
	}
	return ids
}
