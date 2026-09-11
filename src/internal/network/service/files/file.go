package files

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/pkg/filedx"

	"gorm.io/gorm"
)

type FileService struct{}

// Uploads 根据配置的 oss 类型上传文件，返回入库后的文件记录
func (s *FileService) Uploads(fileHeader *multipart.FileHeader) (*basic.Files, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	const maxSize = 32 << 20 // 32MB
	if len(data) > maxSize {
		return nil, fmt.Errorf("文件过大，最大允许 %d MB", maxSize>>20)
	}

	mime, err := filedx.GetFileMIME(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("检测文件类型失败: %w", err)
	}
	if !filedx.CheckIfAllowed(mime.MIME.Value) {
		return nil, fmt.Errorf("不允许的文件类型: %s", mime.MIME.Value)
	}

	hashBytes := md5.Sum(data)
	hashStr := hex.EncodeToString(hashBytes[:])
	ext := "." + mime.Extension

	var existing basic.Files
	if err := global.DB.Where("hash = ?", hashStr).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询文件记录失败: %w", err)
	}

	var url string
	switch global.Config.System.OssType {
	case "local":
		url, err = filedx.UploadToLocal(data, ext)
	case "minio":
		url, err = filedx.UploadToMinio(data, ext)
	case "ali":
		err = fmt.Errorf("阿里云 OSS 上传暂未实现")
	default:
		err = fmt.Errorf("不支持的 oss 类型: %s", global.Config.System.OssType)
	}
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	fileRecord := &basic.Files{
		Url:              url,
		Ext:              ext,
		Size:             int64(len(data)),
		OriginalFileName: fileHeader.Filename,
		Hash:             hashStr,
	}
	if err := global.DB.Create(fileRecord).Error; err != nil {
		return nil, fmt.Errorf("文件记录入库失败: %w", err)
	}
	return fileRecord, nil
}
