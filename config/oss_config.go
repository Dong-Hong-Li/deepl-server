package config

import (
	"errors"
	"os"
	"strconv"
)

type OSSConfig struct {
	// oss 端点
	OSSEndpoint string
	// oss 访问密钥
	OSSAccessKey string
	// oss 密钥
	OSSSecretKey string
	// oss 桶
	OSSBucket string
	// oss 区域 默认 us-east-1
	OSSRegion string
	// oss 预签名 URL 有效期（秒），仅对象存储实现支持时生效。
	OSSPresignGetExpiresSec int
	// 最大文件大小
	MaxFileSize int
}

// 验证 storage 配置
func validateOSSConfig() (*OSSConfig, error) {
	// storage 配置
	ossEndpoint := os.Getenv("APP_OSS_ENDPOINT")
	if ossEndpoint == "" {
		return nil, errors.New("ossEndpoint is required")
	}
	// oss 桶
	ossBucket := os.Getenv("APP_OSS_BUCKET")
	if ossBucket == "" {
		return nil, errors.New("ossBucket is required")
	}

	// oss 访问密钥
	ossAccessKey := os.Getenv("APP_OSS_ACCESS_KEY")
	if ossAccessKey == "" {
		return nil, errors.New("ossAccessKey is required")
	}
	// oss 密钥
	ossSecretKey := os.Getenv("APP_OSS_SECRET_KEY")
	if ossSecretKey == "" {
		return nil, errors.New("ossSecretKey is required")
	}
	// oss 区域
	ossRegion := os.Getenv("APP_OSS_REGION")
	if ossRegion == "" {
		return nil, errors.New("ossRegion is required")
	}
	// storage 预签名 URL 有效期（秒），仅对象存储实现支持时生效。
	ossPresignGetExpiresSecInt, err := strconv.Atoi(os.Getenv("APP_OSS_PRESIGN_GET_EXPIRES_SEC"))
	if err != nil || ossPresignGetExpiresSecInt <= 0 {
		return nil, errors.New("ossPresignGetExpiresSec is invalid")
	}
	// 最大文件大小
	maxFileSize, err := strconv.Atoi(os.Getenv("APP_MAX_FILE_SIZE"))
	if err != nil || maxFileSize <= 0 {
		return nil, errors.New("maxFileSize is invalid")
	}
	return &OSSConfig{
		OSSEndpoint:             ossEndpoint,
		OSSAccessKey:            ossAccessKey,
		OSSSecretKey:            ossSecretKey,
		OSSBucket:               ossBucket,
		OSSRegion:               ossRegion,
		OSSPresignGetExpiresSec: ossPresignGetExpiresSecInt,
		MaxFileSize:             maxFileSize,
	}, nil
}

// StorageConfigured 是否具备对象存储必填项。
func (c *OSSConfig) OSSConfigured() bool {
	return c.OSSEndpoint != "" && c.OSSAccessKey != "" && c.OSSSecretKey != "" && c.OSSBucket != "" && c.OSSRegion != "" && c.OSSPresignGetExpiresSec > 0
}
