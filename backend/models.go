package backend

// ConnectionInfo HDFS连接配置
type ConnectionInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"` // namenode地址，如 "namenode:9000"
	User     string `json:"user"`    // HDFS用户名
	IsActive bool   `json:"isActive"`
}

// FileItem HDFS文件/目录信息
type FileItem struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDir       bool   `json:"isDir"`
	Size        int64  `json:"size"`
	ModTime     string `json:"modTime"`
	Permission  string `json:"permission"`
	Owner       string `json:"owner"`
	Group       string `json:"group"`
	Replication int    `json:"replication"`
}

// ListResult 文件列表结果
type ListResult struct {
	Path  string     `json:"path"`
	Files []FileItem `json:"files"`
	Error string     `json:"error,omitempty"`
}

// OperationResult 操作结果
type OperationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
