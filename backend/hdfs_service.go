package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// HdfsService HDFS文件操作服务（via WebHDFS REST API）
type HdfsService struct {
	ctx     context.Context
	mu      sync.Mutex
	baseURL string // e.g. http://namenode:9879/webhdfs/v1
	user    string
	http    *http.Client
}

func NewHdfsService() *HdfsService {
	return &HdfsService{
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

func (h *HdfsService) Startup(ctx context.Context) {
	h.ctx = ctx
}

// normalizeAddress 将用户输入的地址规范化为 WebHDFS base URL
func normalizeAddress(address string) string {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	return strings.TrimRight(address, "/") + "/webhdfs/v1"
}

// webhdfsURL 构建 WebHDFS 请求 URL，附加 op 和可选的键值对参数
func (h *HdfsService) webhdfsURL(hdfsPath, op string, extraParams ...string) string {
	if !strings.HasPrefix(hdfsPath, "/") {
		hdfsPath = "/" + hdfsPath
	}
	u := h.baseURL + hdfsPath + "?op=" + op
	if h.user != "" {
		u += "&user.name=" + url.QueryEscape(h.user)
	}
	for i := 0; i+1 < len(extraParams); i += 2 {
		u += "&" + extraParams[i] + "=" + url.QueryEscape(extraParams[i+1])
	}
	return u
}

// remoteException 解析 WebHDFS 错误响应
func remoteException(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	var errResp struct {
		RemoteException struct {
			Exception string `json:"exception"`
			Message   string `json:"message"`
		} `json:"RemoteException"`
	}
	if json.NewDecoder(resp.Body).Decode(&errResp) == nil && errResp.RemoteException.Message != "" {
		return fmt.Errorf("%s: %s", errResp.RemoteException.Exception, errResp.RemoteException.Message)
	}
	return fmt.Errorf("HTTP %d %s", resp.StatusCode, resp.Status)
}

// Connect 连接到HDFS（WebHDFS）
func (h *HdfsService) Connect(address, user string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	baseURL := normalizeAddress(address)
	testURL := baseURL + "/?op=LISTSTATUS"
	if user != "" {
		testURL += "&user.name=" + url.QueryEscape(user)
	}

	testClient := &http.Client{Timeout: 15 * time.Second}
	resp, err := testClient.Get(testURL)
	if err != nil {
		hint := ""
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no connection") {
			hint = "\n提示：请确认 WebHDFS 已启用（dfs.webhdfs.enabled=true），端口格式示例：bsa1003:9879"
		}
		return OperationResult{Success: false, Message: "连接失败: " + err.Error() + hint}
	}
	defer resp.Body.Close()

	if err := remoteException(resp); err != nil {
		return OperationResult{Success: false, Message: "连接测试失败: " + err.Error()}
	}

	h.baseURL = baseURL
	h.user = user
	h.http = &http.Client{Timeout: 30 * time.Second}
	return OperationResult{Success: true, Message: "连接成功（WebHDFS）"}
}

// Disconnect 断开连接
func (h *HdfsService) Disconnect() OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.baseURL = ""
	h.user = ""
	return OperationResult{Success: true, Message: "已断开连接"}
}

// IsConnected 检查连接状态
func (h *HdfsService) IsConnected() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.baseURL != ""
}

// webHDFSFileStat WebHDFS 文件状态结构
type webHDFSFileStat struct {
	PathSuffix       string `json:"pathSuffix"`
	Type             string `json:"type"`
	Length           int64  `json:"length"`
	Owner            string `json:"owner"`
	Group            string `json:"group"`
	Permission       string `json:"permission"`
	ModificationTime int64  `json:"modificationTime"`
	Replication      int    `json:"replication"`
}

func statToFileItem(stat webHDFSFileStat, parentPath string) FileItem {
	filePath := path.Join(parentPath, stat.PathSuffix)
	if stat.PathSuffix == "" {
		filePath = parentPath
	}
	return FileItem{
		Name:        stat.PathSuffix,
		Path:        filePath,
		IsDir:       stat.Type == "DIRECTORY",
		Size:        stat.Length,
		ModTime:     time.Unix(stat.ModificationTime/1000, 0).Format(time.DateTime),
		Permission:  stat.Permission,
		Owner:       stat.Owner,
		Group:       stat.Group,
		Replication: stat.Replication,
	}
}

// ListDir 列出目录内容
func (h *HdfsService) ListDir(dirPath string) ListResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.baseURL == "" {
		return ListResult{Path: dirPath, Files: []FileItem{}, Error: "未连接到HDFS"}
	}
	if dirPath == "" {
		dirPath = "/"
	}

	resp, err := h.http.Get(h.webhdfsURL(dirPath, "LISTSTATUS"))
	if err != nil {
		return ListResult{Path: dirPath, Files: []FileItem{}, Error: "读取目录失败: " + err.Error()}
	}
	defer resp.Body.Close()

	if err := remoteException(resp); err != nil {
		return ListResult{Path: dirPath, Files: []FileItem{}, Error: "读取目录失败: " + err.Error()}
	}

	var result struct {
		FileStatuses struct {
			FileStatus []webHDFSFileStat `json:"FileStatus"`
		} `json:"FileStatuses"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ListResult{Path: dirPath, Files: []FileItem{}, Error: "解析响应失败: " + err.Error()}
	}

	files := make([]FileItem, 0, len(result.FileStatuses.FileStatus))
	for _, stat := range result.FileStatuses.FileStatus {
		files = append(files, statToFileItem(stat, dirPath))
	}
	return ListResult{Path: dirPath, Files: files}
}

// GetFileInfo 获取文件信息
func (h *HdfsService) GetFileInfo(filePath string) *FileItem {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.baseURL == "" {
		return nil
	}

	resp, err := h.http.Get(h.webhdfsURL(filePath, "GETFILESTATUS"))
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return nil
	}
	defer resp.Body.Close()

	var result struct {
		FileStatus webHDFSFileStat `json:"FileStatus"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil
	}
	item := statToFileItem(result.FileStatus, path.Dir(filePath))
	item.Path = filePath
	if result.FileStatus.PathSuffix == "" {
		item.Name = path.Base(filePath)
	}
	return &item
}

// MkDir 创建目录
func (h *HdfsService) MkDir(dirPath string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.baseURL == "" {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}

	req, _ := http.NewRequest(http.MethodPut, h.webhdfsURL(dirPath, "MKDIRS", "permission", "755"), nil)
	resp, err := h.http.Do(req)
	if err != nil {
		return OperationResult{Success: false, Message: "创建目录失败: " + err.Error()}
	}
	defer resp.Body.Close()
	if err := remoteException(resp); err != nil {
		return OperationResult{Success: false, Message: "创建目录失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: "目录已创建"}
}

// Delete 删除文件或目录
func (h *HdfsService) Delete(targetPath string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.baseURL == "" {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}

	req, _ := http.NewRequest(http.MethodDelete, h.webhdfsURL(targetPath, "DELETE", "recursive", "true"), nil)
	resp, err := h.http.Do(req)
	if err != nil {
		return OperationResult{Success: false, Message: "删除失败: " + err.Error()}
	}
	defer resp.Body.Close()
	if err := remoteException(resp); err != nil {
		return OperationResult{Success: false, Message: "删除失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: "删除成功"}
}

// Rename 重命名文件或目录
func (h *HdfsService) Rename(oldPath, newPath string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.baseURL == "" {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}

	req, _ := http.NewRequest(http.MethodPut, h.webhdfsURL(oldPath, "RENAME", "destination", newPath), nil)
	resp, err := h.http.Do(req)
	if err != nil {
		return OperationResult{Success: false, Message: "重命名失败: " + err.Error()}
	}
	defer resp.Body.Close()
	if err := remoteException(resp); err != nil {
		return OperationResult{Success: false, Message: "重命名失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: "重命名成功"}
}

// Upload 上传本地文件到HDFS（通过原生对话框选择文件）
func (h *HdfsService) Upload(hdfsDir string) OperationResult {
	localPath, err := runtime.OpenFileDialog(h.ctx, runtime.OpenDialogOptions{
		Title: "选择要上传的文件",
	})
	if err != nil {
		return OperationResult{Success: false, Message: "对话框错误: " + err.Error()}
	}
	if localPath == "" {
		return OperationResult{Success: false, Message: "未选择文件"}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.baseURL == "" {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}

	localFile, err := os.Open(localPath)
	if err != nil {
		return OperationResult{Success: false, Message: "打开本地文件失败: " + err.Error()}
	}
	defer localFile.Close()

	localInfo, err := localFile.Stat()
	if err != nil {
		return OperationResult{Success: false, Message: "获取文件信息失败: " + err.Error()}
	}

	hdfsPath := path.Join(hdfsDir, localInfo.Name())

	// WebHDFS CREATE 是两步 PUT：
	// 1. PUT 到 NameNode → 307 重定向到 DataNode
	// 2. 将文件内容 PUT 到重定向的 DataNode URL
	// 使用不自动跟随重定向的客户端，手动处理
	noRedirectClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 15 * time.Second,
	}

	// Step 1: 触发 NameNode 分配 DataNode 并获取重定向 URL
	req1, _ := http.NewRequest(http.MethodPut, h.webhdfsURL(hdfsPath, "CREATE", "overwrite", "true"), nil)
	resp1, err := noRedirectClient.Do(req1)
	if err != nil {
		return OperationResult{Success: false, Message: "上传初始化失败: " + err.Error()}
	}
	resp1.Body.Close()

	if resp1.StatusCode != http.StatusTemporaryRedirect && resp1.StatusCode != http.StatusCreated {
		return OperationResult{Success: false, Message: fmt.Sprintf("上传初始化异常（HTTP %d）", resp1.StatusCode)}
	}

	dataNodeURL := resp1.Header.Get("Location")
	if dataNodeURL == "" {
		return OperationResult{Success: false, Message: "NameNode 未返回 DataNode 地址"}
	}

	// Step 2: PUT 文件内容到 DataNode
	uploadClient := &http.Client{Timeout: 10 * time.Minute}
	req2, _ := http.NewRequest(http.MethodPut, dataNodeURL, localFile)
	req2.ContentLength = localInfo.Size()
	req2.Header.Set("Content-Type", "application/octet-stream")
	resp2, err := uploadClient.Do(req2)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: "上传到 DataNode 失败: " + err.Error() + "\n提示：DataNode HTTP 端口（通常 9864）可能不可达，请联系管理员",
		}
	}
	defer resp2.Body.Close()
	if err := remoteException(resp2); err != nil {
		return OperationResult{Success: false, Message: "上传失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: fmt.Sprintf("文件 %s 上传成功", localInfo.Name())}
}

// Download 从HDFS下载文件到本地（通过原生对话框选择保存位置）
func (h *HdfsService) Download(hdfsPath string) OperationResult {
	_, fileName := path.Split(hdfsPath)

	savePath, err := runtime.SaveFileDialog(h.ctx, runtime.SaveDialogOptions{
		Title:           "选择保存位置",
		DefaultFilename: fileName,
	})
	if err != nil {
		return OperationResult{Success: false, Message: "对话框错误: " + err.Error()}
	}
	if savePath == "" {
		return OperationResult{Success: false, Message: "未选择保存位置"}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.baseURL == "" {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}

	// WebHDFS OPEN 会重定向到 DataNode，http.Client 默认跟随 GET 重定向
	downloadClient := &http.Client{Timeout: 10 * time.Minute}
	resp, err := downloadClient.Get(h.webhdfsURL(hdfsPath, "OPEN"))
	if err != nil {
		msg := "下载失败: " + err.Error()
		if strings.Contains(err.Error(), "dial tcp") {
			msg += "\n提示：DataNode HTTP 端口（通常 9864）不可达，请联系管理员启用 HttpFS 或配置 NameNode 代理"
		}
		return OperationResult{Success: false, Message: msg}
	}
	defer resp.Body.Close()

	if err := remoteException(resp); err != nil {
		return OperationResult{Success: false, Message: "下载失败: " + err.Error()}
	}

	localFile, err := os.Create(savePath)
	if err != nil {
		return OperationResult{Success: false, Message: "创建本地文件失败: " + err.Error()}
	}
	defer localFile.Close()

	if _, err = io.Copy(localFile, resp.Body); err != nil {
		return OperationResult{Success: false, Message: "写入文件失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: fmt.Sprintf("文件已下载到 %s", savePath)}
}

// GetDiskUsage 获取路径的磁盘使用情况
func (h *HdfsService) GetDiskUsage(dirPath string) map[string]interface{} {
	h.mu.Lock()
	defer h.mu.Unlock()

	result := map[string]interface{}{"success": false}
	if h.baseURL == "" {
		result["message"] = "未连接到HDFS"
		return result
	}

	resp, err := h.http.Get(h.webhdfsURL(dirPath, "GETCONTENTSUMMARY"))
	if err != nil {
		result["message"] = "获取磁盘使用情况失败: " + err.Error()
		return result
	}
	defer resp.Body.Close()

	if err := remoteException(resp); err != nil {
		result["message"] = "获取磁盘使用情况失败: " + err.Error()
		return result
	}

	var cs struct {
		ContentSummary struct {
			Length         int64 `json:"length"`
			FileCount      int64 `json:"fileCount"`
			DirectoryCount int64 `json:"directoryCount"`
		} `json:"ContentSummary"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		result["message"] = "解析响应失败: " + err.Error()
		return result
	}

	result["success"] = true
	result["size"] = cs.ContentSummary.Length
	result["fileCount"] = cs.ContentSummary.FileCount
	result["dirCount"] = cs.ContentSummary.DirectoryCount
	return result
}
