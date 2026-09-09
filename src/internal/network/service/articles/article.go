package articles

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/pkg/commons"
	"errors"
	"go.ube
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleService struct{}

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
		return errors.New("文章发表失败，请稍后重试")
	}
	return nil
}

func (a *ArticleService) Update(req request.UpdateArticleRequest, userId uint) error {
	var article basic.Article
	err := global.DB.Where("id = ?", req.ID).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("更新文章失败：文章不存在", zap.Uint("articleID", req.ID))
			return errors.New("文章不存在")
		}
		global.Log.Error("查询文章失败", zap.Error(err), zap.Uint("articleID", req.ID))
		return errors.New("查询文章失败，请稍后重试")
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
		return errors.New("文章更新失败，请稍后重试")
	}
	if result.RowsAffected == 0 {
		return errors.New("文章更新失败，可能没有变化")
	}
	return nil
}

func (a *ArticleService) Delete(id uint, userId uint) error {
	var article basic.Article
	err := global.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("删除文章失败：文章不存在", zap.Uint("articleID", id))
			return errors.New("文章不存在")
		}
		global.Log.Error("查询文章失败", zap.Error(err), zap.Uint("articleID", id))
		return errors.New("查询文章失败，请稍后重试")
	}

	if err := commons.CheckOwnership(article, userId); err != nil {
		global.Log.Warn("删除文章失败：非文章作者", zap.Uint("articleID", id), zap.Uint("userID", userId))
		return err
	}

	result := global.DB.Delete(&article)
	if result.Error != nil {
		global.Log.Error("删除文章失败", zap.Error(result.Error), zap.Uint("articleID", id))
		return errors.New("文章删除失败，请稍后重试")
	}
	if result.RowsAffected == 0 {
		global.Log.Warn("删除文章失败：文章不存在", zap.Uint("articleID", id))
		return errors.New("文章不存在")
	}
	global.Log.Info("删除文章成功", zap.Uint("articleID", id))
	return nil
}

func (a *ArticleService) GetInfo(id uint) (basic.Article, error) {
	var article basic.Article
	err := global.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return basic.Article{}, errors.New("文章不存在")
		}
		global.Log.Error("查询文章详情失败", zap.Error(err), zap.Uint("articleID", id))
		return basic.Article{}, errors.New("查询文章失败，请稍后重试")
	}
	return article, nil
}

func (a *ArticleService) ListByUser(userId uint, pageReq request.PageReq) ([]basic.Article, int64, error) {
	query := global.DB.Model(&basic.Article{}).Where("user_id = ?", userId)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.Log.Error("统计用户文章总数失败", zap.Error(err), zap.Uint("userID", userId))
		return nil, 0, errors.New("查询用户文章列表失败，请稍后重试")
	}

	var articles []basic.Article
	if err := query.Order("is_top DESC, id DESC").
		Offset(pageReq.Offset()).
		Limit(pageReq.PageSize).
		Find(&articles).Error; err != nil {
		global.Log.Error("查询用户文章列表失败", zap.Error(err), zap.Uint("userID", userId))
		return nil, 0, errors.New("查询用户文章列表失败，请稍后重试")
	}

	return articles, total, nil
}

func (a *ArticleService) List(req request.ArticleQueryListReq) ([]basic.Article, int64, error) {
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
		return nil, 0, errors.New("查询文章列表失败，请稍后重试")
	}

	var articles []basic.Article
	if err := query.Order("is_top DESC, id DESC").
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&articles).Error; err != nil {
		global.Log.Error("查询文章列表失败", zap.Error(err))
		return nil, 0, errors.New("查询文章列表失败，请稍后重试")
	}

	return articles, total, nil
}
