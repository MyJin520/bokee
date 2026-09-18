package articles

import (
	"bokee/internal/network/api/articles"
	"bokee/pkg/routex"
	"github.com/gin-gonic/gin"
)

var articleApi = articles.ArticleApi{}

func InitArticleRouter(privateGroup *gin.RouterGroup) {
	articlePrivate := routex.NewGroup("文章模块", privateGroup.Group("/article"))
	{
		articlePrivate.POST("/create", "创建文章", articleApi.Create)                  // 创建文章
		articlePrivate.PUT("/update", "更新文章", articleApi.Update)                   // 更新文章
		articlePrivate.DELETE("/delete", "删除文章", articleApi.Delete)                // 删除文章
		articlePrivate.GET("/get_info", "获取文章详情", articleApi.GetInfo)              // 获取文章详情
		articlePrivate.POST("/list", "分页获取文章列表", articleApi.List)                  // 分页获取文章列表
		articlePrivate.POST("/list_by_user", "查看指定用户的所有文章", articleApi.ListByUser) // 查看指定用户的所有文章
	}
}
