package backend

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/colinmarc/hdfs/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// HdfsService HDFS文件操作服务
type HdfsService struct {
	ctx    context.Context
	mu     sync.Mutex
	client *hdfs.Client
}

func NewHdfsService() *HdfsService {
	return &HdfsService{}
}

func (h *HdfsService) Startup(ctx context.Context) {
	h.ctx = ctx
}

// Connect 连接到HDFS
func (h *HdfsService) Connect(address, user string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 先断开已有连接
	if h.client != nil {
		h.client.Close()
		h.client = nil
	}

	opts := hdfs.ClientOptions{
		Addresses: []string{address},
	}
	if user != "" {
		opts.User = user
	}

	client, err := hdfs.NewClient(opts)
	if err != nil {
		return OperationResult{Success: false, Message: "连接失败: " + err.Error()}
	}

	// 测试连接
	_, err = client.ReadDir("/")
	if err != nil {
		client.Close()
		return OperationResult{Success: false, Message: "连接测试失败: " + err.Error()}
	}

	h.client = client
	return OperationResult{Success: true, Message: "连接成功"}
}

// Disconnect 断开连接
func (h *HdfsService) Disconnect() OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client != nil {
		h.client.Close()
		h.client = nil
	}
	return OperationResult{Success: true, Message: "已断开连接"}
}

// IsConnected 检查连接状态
func (h *HdfsService) IsConnected() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.client != nil
}

// ListDir 列出目录内容
func (h *HdfsService) ListDir(dirPath string) ListResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client == nil {
		return ListResult{Path: dirPath, Files: []FileItem{}, Error: "未连接到HDFS"}
	}

	if dirPath == "" {
		dirPath = "/"
	}

	entries, err := h.client.ReadDir(dirPath)
	if err != nil {
		return ListResult{Path: dirPath, Files: []FileItem{}, Error: "读取目录失败: " + err.Error()}
	}

	var files []FileItem
	for _, entry := range entries {
		fi := FileItem{
			Name:    entry.Name(),
			Path:    path.Join(dirPath, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    entry.Size(),
			ModTime: entry.ModTime().Format(time.DateTime),
		}
		// 获取HDFS特有信息
		if stat, ok := entry.(*hdfs.FileInfo); ok {
			fi.Permission = stat.Mode().Perm().String()
			fi.Owner = stat.Owner()
			fi.Group = stat.OwnerGroup()
			// 目录没有replication
			if !entry.IsDir() {
				// 不直接可用，置为0
			}
		}
		files = append(files, fi)
	}

	if files == nil {
		files = []FileItem{}
	}
	return ListResult{Path: dirPath, Files: files}
}

// GetFileInfo 获取文件信息
func (h *HdfsService) GetFileInfo(filePath string) *FileItem {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client == nil {
		return nil
	}

	info, err := h.client.Stat(filePath)
	if err != nil {
		return nil
	}

	fi := &FileItem{
		Name:    info.Name(),
		Path:    filePath,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime().Format(time.DateTime),
	}
	if stat, ok := info.(*hdfs.FileInfo); ok {
		fi.Permission = stat.Mode().Perm().String()
		fi.Owner = stat.Owner()
		fi.Group = stat.OwnerGroup()
	}
	return fi
}

// MkDir 创建目录
func (h *HdfsService) MkDir(dirPath string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client == nil {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}
	err := h.client.MkdirAll(dirPath, 0755)
	if err != nil {
		return OperationResult{Success: false, Message: "创建目录失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: "目录已创建"}
}

// Delete 删除文件或目录
func (h *HdfsService) Delete(targetPath string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client == nil {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}

	info, err := h.client.Stat(targetPath)
	if err != nil {
		return OperationResult{Success: false, Message: "文件不存在: " + err.Error()}
	}

	if info.IsDir() {
		err = h.client.RemoveAll(targetPath)
	} else {
		err = h.client.Remove(targetPath)
	}

	if err != nil {
		return OperationResult{Success: false, Message: "删除失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: "删除成功"}
}

// Rename 重命名文件或目录
func (h *HdfsService) Rename(oldPath, newPath string) OperationResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client == nil {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}
	err := h.client.Rename(oldPath, newPath)
	if err != nil {
		return OperationResult{Success: false, Message: "重命名失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: "重命名成功"}
}

// Upload 上传本地文件到HDFS（通过原生对话框选择文件）
func (h *HdfsService) Upload(hdfsDir string) OperationResult {
	// 选择文件对话框在锁外执行
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

	if h.client == nil {
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
	writer, err := h.client.Create(hdfsPath)
	if err != nil {
		return OperationResult{Success: false, Message: "创建HDFS文件失败: " + err.Error()}
	}

	_, err = io.Copy(writer, localFile)
	if err != nil {
		writer.Close()
		return OperationResult{Success: false, Message: "上传失败: " + err.Error()}
	}

	if err := writer.Close(); err != nil {
		return OperationResult{Success: false, Message: "关闭HDFS写入流失败: " + err.Error()}
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

	if h.client == nil {
		return OperationResult{Success: false, Message: "未连接到HDFS"}
	}

	reader, err := h.client.Open(hdfsPath)
	if err != nil {
		return OperationResult{Success: false, Message: "打开HDFS文件失败: " + err.Error()}
	}
	defer reader.Close()

	localFile, err := os.Create(savePath)
	if err != nil {
		return OperationResult{Success: false, Message: "创建本地文件失败: " + err.Error()}
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, reader)
	if err != nil {
		return OperationResult{Success: false, Message: "下载失败: " + err.Error()}
	}

	return OperationResult{Success: true, Message: fmt.Sprintf("文件已下载到 %s", savePath)}
}

// GetDiskUsage 获取路径的磁盘使用情况
func (h *HdfsService) GetDiskUsage(dirPath string) map[string]interface{} {
	h.mu.Lock()
	defer h.mu.Unlock()

	result := map[string]interface{}{
		"success": false,
	}

	if h.client == nil {
		result["message"] = "未连接到HDFS"
		return result
	}

	info, err := h.client.GetContentSummary(dirPath)
	if err != nil {
		result["message"] = "获取磁盘使用情况失败: " + err.Error()
		return result
	}

	result["success"] = true
	result["size"] = info.Size()
	result["fileCount"] = info.FileCount()
	result["dirCount"] = info.DirectoryCount()
	return result
}
