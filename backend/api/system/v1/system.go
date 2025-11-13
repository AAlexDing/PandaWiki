package v1

type SystemReq struct {
	KbID string `json:"kb_id" query:"kb_id" validate:"required"`
}

type SystemResp struct {
	Document DocumentInfo `json:"document"`
	Learning LearningInfo `json:"learning"`
	System   SystemInfo   `json:"system"`
}

type DocumentInfo struct {
	CurrentCount      int64 `json:"current_count"`       // 当前文档数
	NewIn24h          int64 `json:"new_in_24h"`          // 24h新增文档数
	LearningSucceeded int64 `json:"learning_succeeded"` // 学习成功数量
	LearningFailed    int64 `json:"learning_failed"`     // 学习失败数量
}

type LearningInfo struct {
	BasicProcessing    QueueProgress `json:"basic_processing"`     // 基础处理队列进度
	BasicFailed       int64         `json:"basic_failed"`          // 基础处理失败数
	EnhanceProcessing QueueProgress `json:"enhance_processing"`   // 增强处理队列进度
	EnhanceFailed     int64         `json:"enhance_failed"`        // 增强处理失败数
	BasicFailedDocs   []FailedDoc   `json:"basic_failed_docs"`    // 基础处理失败文档
	EnhanceFailedDocs []FailedDoc   `json:"enhance_failed_docs"` // 增强处理失败文档
}

type QueueProgress struct {
	Pending  int64 `json:"pending"`  // 等待中
	Running  int64 `json:"running"`  // 运行中
	Total    int64 `json:"total"`    // 总数
	Progress int   `json:"progress"` // 进度百分比 (0-100)
}

type FailedDoc struct {
	NodeID   string `json:"node_id"`   // 节点ID
	NodeName string `json:"node_name"` // 文档名
	Reason   string `json:"reason"`   // 失败原因
}

type SystemInfo struct {
	Components []ComponentStatus `json:"components"`
}

type ComponentStatus struct {
	Name      string `json:"name"`       // 组件名称
	Status    string `json:"status"`     // 状态: running, stopped, error
	Image     string `json:"image"`      // 镜像名称
	Ports     string `json:"ports"`       // 端口信息
	Health    string `json:"health"`     // 健康状态 (仅RAGLite和Qdrant)
	LogStatus string `json:"log_status"` // 日志解析状态 (仅RAGLite和Qdrant)
}
