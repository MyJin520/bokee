package articles

import (
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/internal/network/middleware"
	"gin-admin/internal/network/service/articles"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ArticleApi struct{}

var articleService = &articles.ArticleService{}

func (a ArticleApi) Create(c *gin.Context) {
	var req request.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	userId, _ := middleware.GetUserIDFromContext(c)
	err := articleService.Create(req, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("文章发表成功", c)
}

func (a ArticleApi) Update(c *gin.Context) {
	var req request.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	err := articleService.Update(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("文章更新成功", c)
}

func (a ArticleApi) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("文章ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32) // 转为 uint32 或 uint
	if err != nil {
		response.FailWithMessage("文章ID格式错误", c)
		return
	}

	articleID := uint(id)
	if err := articleService.Delete(articleID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("文章删除成功", c)
}

func (a ArticleApi) GetInfo(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("文章ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32) // 转为 uint32 或 uint
	if err != nil {
		response.FailWithMessage("文章ID格式错误", c)
		return
	}
	articleID := uint(id)
	article, err := articleService.GetInfo(articleID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(article, "文章详情获取成功", c)
}

func (a ArticleApi) ListByUser(c *gin.Context) {
	var req request.UserArticleListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	req.Normalize()

	articleList, total, err := articleService.ListByUser(req.UserID, req.PageReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(articleList, total, req.Page, req.PageSize, "用户文章列表获取成功", c)
}

func (a ArticleApi) List(c *gin.Context) {
	var req request.ArticleQueryListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	articleList, total, err := articleService.List(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(articleList, total, req.Page, req.PageSize, "文章列表获取成功", c)
}
