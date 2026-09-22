package articles

import (
	"context"
	"fmt"
	"time"

	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/pkg/commons"
	"bokee/pkg/redisx"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleService struct{}

// articleInfoKey 文章详情缓存 Key：bokee:article:info:<id>
func articleInfoKey(id uint) string {
	return redisx.BuildKey("article", "info", fmt.Sprintf("%d", id))
}

// deleteArticleInfoCache 写操作落库成功后失效文章详情缓存（先写库、后删缓存，保证最终一致）
func deleteArticleInfoCache(ctx context.Context, ids ...uint) {
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, articleInfoKey(id))
	}
	if err := redisx.Delete(ctx, keys...); err != nil {
		global.Log.Error("删除文章详情缓存失败", zap.Strings("keys", keys), zap.Error(err))
	}
}

// DeleteArticleInfoCacheByUser 失效指定用户的所有文章详情缓存（如用户更新姓名/头像时调用）
func (a *ArticleService) DeleteArticleInfoCacheByUser(ctx context.Context, userId uint) {
	var ids []uint
	if err := global.DB.Model(&basic.Article{}).Where("user_id = ?", userId).Pluck("id", &ids).Error; err != nil {
		global.Log.Warn("查询用户文章ID失败，跳过缓存失效", zap.Error(err), zap.Uint("userID", userId))
		return
	}
	if len(ids) == 0 {
		return
	}
	deleteArticleInfoCache(ctx, ids...)
}

func (a *ArticleService) Create(req request.CreateArticleRequest, userId uint) error {
	newArticle := basic.Article{
		UserID:  userId,
		Title:   req.Title,
		Content: req.Content,
		Summary: req.Summary,
		Cover:   req.Cover,
	}
	if err := global.DB.Create(&newArticle).Error; err != nil {
		global.Log.Error("发表文章失败", zap.Error(err), zap.Uint("userID", userId))
		return fmt.Errorf("文章发表失败，请稍后重试")
	}
	return nil
}

func (a *ArticleService) Update(ctx context.Context, req request.UpdateArticleRequest, userId uint) error {
	var article basic.Article
	err := global.DB.Where("id = ?", req.ID).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("更新文章失败：文章不存在", zap.Uint("articleID", req.ID))
			return fmt.Errorf("文章不存在")
		}
		global.Log.Error("查询文章失败", zap.Error(err), zap.Uint("articleID", req.ID))
		return fmt.Errorf("查询文章失败，请稍后重试")
	}

	if err := commons.CheckOwnership(article, userId); err != nil {
		global.Log.Warn("更新文章失败：非文章作者", zap.Uint("articleID", req.ID), zap.Uint("userID", userId))
		return err
	}

	updates := commons.StructToUpdateMap(req)
	delete(updates, "id")

	if len(updates) == 0 {
		return nil
	}

	result := global.DB.Model(&article).Updates(updates)
	if result.Error != nil {
		global.Log.Error("更新文章失败", zap.Error(result.Error), zap.Uint("articleID", req.ID))
		return fmt.Errorf("文章更新失败，请稍后重试")
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("文章更新失败，可能没有变化")
	}

	// 落库成功后失效缓存
	deleteArticleInfoCache(ctx, article.ID)
	return nil
}

func (a *ArticleService) Delete(ctx context.Context, id uint, userId uint) error {
	var article basic.Article
	err := global.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("删除文章失败：文章不存在", zap.Uint("articleID", id))
			return fmt.Errorf("文章不存在")
		}
		global.Log.Error("查询文章失败", zap.Error(err), zap.Uint("articleID", id))
		return fmt.Errorf("查询文章失败，请稍后重试")
	}

	if err := commons.CheckOwnership(article, userId); err != nil {
		global.Log.Warn("删除文章失败：非文章作者", zap.Uint("articleID", id), zap.Uint("userID", userId))
		return err
	}

	result := global.DB.Delete(&article)
	if result.Error != nil {
		global.Log.Error("删除文章失败", zap.Error(result.Error), zap.Uint("articleID", id))
		return fmt.Errorf("文章删除失败，请稍后重试")
	}
	if result.RowsAffected == 0 {
		global.Log.Warn("删除文章失败：文章不存在", zap.Uint("articleID", id))
		return fmt.Errorf("文章不存在")
	}
	global.Log.Info("删除文章成功", zap.Uint("articleID", id))
	// 落库成功后失效缓存
	deleteArticleInfoCache(ctx, article.ID)
	return nil
}

// buildArticleInfoResp 将文章模型转换为详情响应结构体
func buildArticleInfoResp(article basic.Article, user basic.User) response.ArticleInfoResp {
	return response.ArticleInfoResp{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Summary:   article.Summary,
		Cover:     article.Cover,
		ViewCount: article.ViewCount,
		LikeCount: article.LikeCount,
		IsTop:     article.IsTop,
		UserID:    article.UserID,
		Author: response.AuthorInfo{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		},
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
}

// buildArticleListItemResp 将文章模型转换为列表响应结构体（不含正文内容）
func buildArticleListItemResp(article basic.Article) response.ArticleListItemResp {
	return response.ArticleListItemResp{
		ID:        article.ID,
		Title:     article.Title,
		Summary:   article.Summary,
		Cover:     article.Cover,
		ViewCount: article.ViewCount,
		LikeCount: article.LikeCount,
		IsTop:     article.IsTop,
		UserID:    article.UserID,
		CreatedAt: article.CreatedAt,
	}
}

func (a *ArticleService) GetInfo(ctx context.Context, id uint) (response.ArticleInfoResp, error) {
	// 先查缓存
	key := articleInfoKey(id)
	var cached response.ArticleInfoResp
	if err := redisx.GetJSON(ctx, key, &cached); err != nil {
		global.Log.Error("读取文章详情缓存失败，降级查库", zap.Error(err), zap.Uint("articleID", id))
	} else if cached.ID != 0 {
		return cached, nil
	}

	// 缓存 miss，查库
	var article basic.Article
	err := global.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ArticleInfoResp{}, fmt.Errorf("文章不存在")
		}
		global.Log.Error("查询文章详情失败", zap.Error(err), zap.Uint("articleID", id))
		return response.ArticleInfoResp{}, fmt.Errorf("查询文章失败，请稍后重试")
	}

	// 关联查询作者信息
	var user basic.User
	if err := global.DB.Select("id", "name", "avatar").Where("id = ?", article.UserID).First(&user).Error; err != nil {
		global.Log.Warn("查询文章作者信息失败", zap.Error(err), zap.Uint("userID", article.UserID))
	}

	// 查库成功后回填缓存
	resp := buildArticleInfoResp(article, user)
	if err := redisx.SetJSON(ctx, key, resp, 30*time.Minute); err != nil {
		global.Log.Error("回填文章详情缓存失败", zap.Error(err), zap.Uint("articleID", id))
	}
	return resp, nil
}

func (a *ArticleService) ListByUser(userId uint, pageReq request.PageReq) ([]response.ArticleListItemResp, int64, error) {
	query := global.DB.Model(&basic.Article{}).Where("user_id = ?", userId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Log.Error("统计用户文章总数失败", zap.Error(err), zap.Uint("userID", userId))
		return nil, 0, fmt.Errorf("查询用户文章列表失败，请稍后重试")
	}

	var articles []basic.Article
	if err := query.Order("is_top DESC, id DESC").
		Offset(pageReq.Offset()).
		Limit(pageReq.PageSize).
		Find(&articles).Error; err != nil {
		global.Log.Error("查询用户文章列表失败", zap.Error(err), zap.Uint("userID", userId))
		return nil, 0, fmt.Errorf("查询用户文章列表失败，请稍后重试")
	}

	list := make([]response.ArticleListItemResp, 0, len(articles))
	for _, article := range articles {
		list = append(list, buildArticleListItemResp(article))
	}
	return list, total, nil
}

func (a *ArticleService) List(req request.ArticleQueryListReq) ([]response.ArticleListItemResp, int64, error) {
	query := global.DB.Model(&basic.Article{})

	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Summary != "" {
		query = query.Where("summary LIKE ?", "%"+req.Summary+"%")
	}
	if req.UserName != "" {
		query = query.Joins("JOIN users ON users.id = articles.user_id").
			Where("users.name LIKE ?", "%"+req.UserName+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Log.Error("统计文章总数失败", zap.Error(err))
		return nil, 0, fmt.Errorf("查询文章列表失败，请稍后重试")
	}

	var articles []basic.Article
	if err := query.Order("is_top DESC, id DESC").
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&articles).Error; err != nil {
		global.Log.Error("查询文章列表失败", zap.Error(err))
		return nil, 0, fmt.Errorf("查询文章列表失败，请稍后重试")
	}

	list := make([]response.ArticleListItemResp, 0, len(articles))
	for _, article := range articles {
		list = append(list, buildArticleListItemResp(article))
	}
	return list, total, nil
}