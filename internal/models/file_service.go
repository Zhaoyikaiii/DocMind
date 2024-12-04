package models

// FileListParams 定义文件列表查询参数
type FileListParams struct {
	UploaderID  *uint   // 上传者ID
	DocumentID  *uint   // 关联的文档ID
	ContentType *string // 文件类型
	Search      string  // 搜索关键词
	Page        int     // 页码
	PageSize    int     // 每页数量
}
