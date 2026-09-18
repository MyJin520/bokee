package files

import (
	"bokee/internal/network/api/files"
	"bokee/pkg/routex"
	"github.com/gin-gonic/gin"
)

var fileApi = files.FileApi{}

func InitFileRouter(publicGroup *gin.RouterGroup) {
	filePublic := routex.NewGroup("文件模块", publicGroup.Group("/file"))
	{
		filePublic.POST("/uploads", "文件上传", fileApi.Uploads) // 文件上传
	}
}
