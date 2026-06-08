package articles

import (
	"gin-admin/internal/network/api/articles"
	"github.com/gin-gonic/gin"
)

var articleApi = articles.ArticleApi{}

func InitArticleRouter(privateGroup *gin.RouterGroup) {
	articlePrivate := privateGroup.Group("/article")
	{
		articlePrivate.POST("/create", articleApi.Create)   // 创建文章
		articlePrivate.PUT("/update", articleApi.Update)    // 更新文章
		articlePrivate.DELETE("/delete", articleApi.Delete) // 删除文章
		articlePrivate.GET("/get_info", articleApi.GetInfo) // 获取文章详情
		articlePrivate.POST("/list", articleApi.List)       // 分页获取文章列表
	}
}
