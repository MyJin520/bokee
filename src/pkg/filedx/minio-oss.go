package filedx

import (
	"bokee/global"
	"bokee/pkg/randx"
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// initMinioClient 初始化 MinIO 客户端，并确保目标桶存在。
func initMinioClient() (*minio.Client, error) {
	cfg := global.Config.Oss.MinioOss
	client, err := minio.New(
		cfg.Host,
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
			Secure: cfg.UseSSL,
		},
	)
	if err != nil {
		return nil, errors.New("minio客户端初始化失败：" + err.Error())
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, errors.New("检查桶失败：" + err.Error())
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, errors.New("创建桶失败：" + err.Error())
		}
	}
	return client, nil
}

// UploadToMinio 将文件字节上传到 MinIO，返回文件的可访问 URL。
func UploadToMinio(data []byte, ext string) (string, error) {
	cfg := global.Config.Oss.MinioOss

	client, err := initMinioClient()
	if err != nil {
		return "", err
	}

	// 按日期分目录，并使用 UUID 保证对象名唯一，与本地存储保持一致。
	objectName := fmt.Sprintf("%s/%s.%s", time.Now().Format("2006/01/02"), randx.RandomUUID(), ext)

	_, err = client.PutObject(
		context.Background(),
		cfg.BucketName,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: mime.TypeByExtension(ext)},
	)
	if err != nil {
		return "", errors.New("上传文件到 MinIO 失败：" + err.Error())
	}

	// 拼接对外访问地址：{file-url}/{bucket}/{objectName}
	base := strings.TrimRight(cfg.FileUrl, "/")
	if base == "" {
		return objectName, nil
	}
	return fmt.Sprintf("%s/%s/%s", base, cfg.BucketName, objectName), nil
}
