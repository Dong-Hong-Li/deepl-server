package objectstorage

import (
	"context"
	"deep-map-server/config"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	Storage *config.OSSConfig
	// S3 客户端
	client *s3.Client
	// 存储桶名称
	bucket string
	// 预签名下载链接过期时间
	presignGetExpires time.Duration
	// 初始化存储桶
	initBucket sync.Once
}

func NewClient(ossConfig *config.OSSConfig, ctx context.Context) (*Client, error) {
	if ossConfig == nil || !ossConfig.OSSConfigured() {
		return nil, errors.New("ossConfig is required")
	}
	// 初始化S3 客户端配置
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(ossConfig.OSSRegion),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(ossConfig.OSSAccessKey, ossConfig.OSSSecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	//创建S3 客户端
	s3service := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(ossConfig.OSSEndpoint)
	})

	// 预签名下载链接过期时间
	sec := ossConfig.OSSPresignGetExpiresSec

	return &Client{Storage: ossConfig,
		client:            s3service,
		bucket:            ossConfig.OSSBucket,
		presignGetExpires: time.Duration(sec) * time.Second,
	}, nil
}

func (c *Client) ensureBucket(ctx context.Context) error {
	var initErr error
	c.initBucket.Do(func() {
		_, err := c.client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(c.bucket),
		})
		if err == nil {
			return
		}

		// 如果存储桶不存在，则创建存储桶
		_, createErr := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(c.bucket),
		})
		if createErr != nil {
			initErr = fmt.Errorf("ensure bucket: %w", createErr)
		}
	})
	return initErr
}

// Upload 上传文件到 S3
func (c *Client) Upload(ctx context.Context, key string, body io.Reader, contentType string) error {
	if err := c.ensureBucket(ctx); err != nil {
		return err
	}

	// 上传文件到 S3
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

// PresignGetObjectURL 生成短期下载链接（SigV4 预签名 GET）。
func (c *Client) PresignGetObjectURL(ctx context.Context, objectKey string) (string, error) {
	if objectKey == "" {
		return "", errors.New("empty object key")
	}
	pc := s3.NewPresignClient(c.client)
	out, err := pc.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = c.presignGetExpires
	})
	if err != nil {
		return "", fmt.Errorf("presign get object: %w", err)
	}
	return out.URL, nil
}

// GetObject 从桶中拉取对象正文；用于重分析时按 storage_key 回灌 resume_text，避免再次解析 PDF。
func (c *Client) GetObject(ctx context.Context, key string) ([]byte, string, error) {
	if key == "" {
		return nil, "", errors.New("empty object key")
	}
	out, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", fmt.Errorf("get object: %w", err)
	}
	defer out.Body.Close()
	data, rerr := io.ReadAll(out.Body)
	if rerr != nil {
		return nil, "", fmt.Errorf("read object body: %w", rerr)
	}
	ct := ""
	if out.ContentType != nil {
		ct = *out.ContentType
	}
	return data, ct, nil
}

// DeleteObject 从桶中删除对象；key 为空时为 no-op。
func (c *Client) DeleteObject(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}
