package articles

import (
	"errors"
	"gin-admin/global"
	"gin-admin/internal/mods/basic"
	"gin-admin/internal/mods/request"
	"gin-admin/pkg/commons"
	"gorm.io/gorm"
)

type ArticleService struct{}

func (s ArticleService) CreateArticle(req request.CreateArticleRequest, userId uint) error {
	newArticle := basic.Article{
		UserID:  userId,
		Title:   req.Title,
		Content: req.Content,
		Summary: req.Summary,
		Cover:   req.Cover,
	}
	err := global.DB.Create(&newArticle).Error
	if err != nil {
		return errors.New("文章发表失败")
	}
	return nil
}

func (s ArticleService) UpdateArticle(req request.UpdateArticleRequest) error {
	var article basic.Article
	err := global.DB.Where("id = ?", req.ID).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return errors.New("查询文章失败: " + err.Error())
	}

	updates := commons.StructToUpdateMap(req)
	delete(updates, "id")

	if len(updates) == 0 {
		return nil
	}

	result := global.DB.Model(&article).Updates(updates)
	if result.Error != nil {
		return errors.New("文章更新失败: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return errors.New("文章更新失败，可能没有变化")
	}
	return nil
}

func (s ArticleService) DeleteArticle(id uint) error {
	result := global.DB.Where("id = ?", id).Delete(&basic.Article{})
	if result.Error != nil {
		return errors.New("文章删除失败: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return errors.New("文章不存在")
	}
	return nil
}

func (s ArticleService) GetArticleInfo(id uint) (basic.Article, error) {
	var article basic.Article
	err := global.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return basic.Article{}, errors.New("文章不存在")
		}
		return basic.Article{}, errors.New("查询文章失败: " + err.Error())
	}
	return article, nil
}

func (s ArticleService) GetArticleList(req request.ArticleQueryListReq) ([]basic.Article, int64, error) {
	return nil, 0, nil
}
