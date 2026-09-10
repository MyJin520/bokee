package files

import (
	"bokee/internal/network/api/files"
	"github.com/gin-gonic/gin"
)

var fileApi = files.FileApi{}

func InitFileRouter(publicGroup *gin.RouterGroup) {
	filePublic := publicGroup.Group("/file")
	{
		filePublic.POST("/uploads", fileApi.Uploads)
	}
}
