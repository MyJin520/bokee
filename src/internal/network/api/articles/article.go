package articles

import (
	"bokee/global"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/internal/network/middleware"
	"bokee/internal/network/service/articles"
	"bokee/pkg/verifyx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type ArticleApi struct{}

var articleService = &articles.ArticleService{}

func (a *ArticleApi) Create(c *gin.Context) {
	var req request.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("文章创建请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("文章创建参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
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

func (a *ArticleApi) Update(c *gin.Context) {
	var req request.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("文章更新请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("文章更新参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}
	userId, _ := middleware.GetUserIDFromContext(c)
	err := articleService.Update(req, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("文章更新成功", c)
}

func (a *ArticleApi) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithRequest("文章ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32) // 转为 uint32 或 uint
	if err != nil {
		global.Log.Warn("文章ID格式错误", zap.String("id", idStr))
		response.FailWithRequest("文章ID格式错误", c)
		return
	}

	articleID := uint(id)
	userId, _ := middleware.GetUserIDFromContext(c)
	if err := articleService.Delete(articleID, userId); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("文章删除成功", c)
}

func (a *ArticleApi) GetInfo(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithRequest("文章ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32) // 转为 uint32 或 uint
	if err != nil {
		global.Log.Warn("文章ID格式错误", zap.String("id", idStr))
		response.FailWithRequest("文章ID格式错误", c)
		return
	}
	articleID := uint(id)
	article, err := articleService.GetInfo(articleID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(article, "文章详情获取成功", c)
}

func (a *ArticleApi) ListByUser(c *gin.Context) {
	var req request.UserArticleListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("用户文章列表请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("用户文章列表参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
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

func (a *ArticleApi) List(c *gin.Context) {
	var req request.ArticleQueryListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("文章列表请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("文章列表参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}
	articleList, total, err := articleService.List(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(articleList, total, req.Page, req.PageSize, "文章列表获取成功", c)
}
