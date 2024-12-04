package storage

import (
	"fmt"

	"github.com/Zhaoyikaiii/docmind/internal/storage/operators"
)

// StorageType 存储类型
type StorageType string

const (
	Local StorageType = "local"
	OSS   StorageType = "oss"
	S3    StorageType = "s3"
	COS   StorageType = "cos"
	Qiniu StorageType = "qiniu"
)

// NewFileOperator 创建文件操作器
func NewFileOperator(storageType StorageType, config interface{}) (FileOperator, error) {
	switch storageType {
	case Local:
		cfg, ok := config.(operators.LocalConfig)
		if !ok {
			return nil, fmt.Errorf("invalid config type for local storage")
		}
		return operators.NewLocalFileOperator(cfg), nil

	case OSS:
		cfg, ok := config.(operators.OSSConfig)
		if !ok {
			return nil, fmt.Errorf("invalid config type for OSS")
		}
		return operators.NewOSSFileOperator(cfg)

	case S3:
		cfg, ok := config.(operators.S3Config)
		if !ok {
			return nil, fmt.Errorf("invalid config type for S3")
		}
		return operators.NewS3FileOperator(cfg)

	case COS:
		cfg, ok := config.(operators.COSConfig)
		if !ok {
			return nil, fmt.Errorf("invalid config type for COS")
		}
		return operators.NewCOSFileOperator(cfg)

	case Qiniu:
		cfg, ok := config.(operators.QiniuConfig)
		if !ok {
			return nil, fmt.Errorf("invalid config type for Qiniu")
		}
		return operators.NewQiniuFileOperator(cfg), nil

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
