package files

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/pkg/filedx"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"gorm.io/gorm"
)

type FileService struct{}

// Uploads 根据配置的 oss 类型上传文件，返回入库后的文件记录
func (s *FileService) Uploads(fileHeader *multipart.FileHeader) (*basic.Files, error) {
	// 1. 读取文件内容并校验类型
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	mime, err := filedx.GetFileMIME(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	if !filedx.CheckIfAllowed(mime.MIME.Value) {
		return nil, fmt.Errorf("不允许的文件类型: %s", mime.MIME.Value)
	}

	// 2. 计算文件 hash 与扩展名
	hashBytes := md5.Sum(data)
	hashStr := hex.EncodeToString(hashBytes[:])
	ext := "." + mime.Extension

	// 3. 根据 hash 去重：已存在相同内容则直接复用记录
	var existing basic.Files
	if err := global.DB.Where("hash = ?", hashStr).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询文件记录失败: %w", err)
	}

	// 4. 根据 oss 类型选择上传位置
	ossType := global.Config.System.OssType
	var url string
	switch ossType {
	case "local":
		url, err = filedx.UploadToLocal(data, ext)
	case "minio":
		url, err = filedx.UploadToMinio(data, ext)
	case "ali":
		err = errors.New("阿里云 OSS 上传暂未实现")
	default:
		err = fmt.Errorf("不支持的 oss 类型: %s", ossType)
	}
	if err != nil {
		return nil, err
	}

	// 5. 入库
	fileRecord := &basic.Files{
		Url:              url,
		Ext:              ext,
		Size:             fileHeader.Size,
		OriginalFileName: fileHeader.Filename,
		Hash:             hashStr,
	}
	if err := global.DB.Create(fileRecord).Error; err != nil {
		return nil, fmt.Errorf("文件记录入库失败: %w", err)
	}
	return fileRecord, nil
}
