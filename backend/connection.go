package backend

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

// ConnectionManager 管理HDFS连接配置
type ConnectionManager struct {
	ctx         context.Context
	mu          sync.Mutex
	connections []ConnectionInfo
	configPath  string
}

func NewConnectionManager() *ConnectionManager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".hdfs-king")
	os.MkdirAll(configDir, 0755)
	cm := &ConnectionManager{
		configPath: filepath.Join(configDir, "connections.json"),
	}
	cm.loadConnections()
	return cm
}

func (cm *ConnectionManager) Startup(ctx context.Context) {
	cm.ctx = ctx
}

// GetConnections 获取所有连接
func (cm *ConnectionManager) GetConnections() []ConnectionInfo {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.connections == nil {
		return []ConnectionInfo{}
	}
	return cm.connections
}

// AddConnection 添加新连接
func (cm *ConnectionManager) AddConnection(name, address, user string) OperationResult {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if name == "" || address == "" {
		return OperationResult{Success: false, Message: "名称和地址不能为空"}
	}

	conn := ConnectionInfo{
		ID:      uuid.New().String(),
		Name:    name,
		Address: address,
		User:    user,
	}
	cm.connections = append(cm.connections, conn)
	if err := cm.saveConnections(); err != nil {
		return OperationResult{Success: false, Message: "保存连接失败: " + err.Error()}
	}
	return OperationResult{Success: true, Message: "连接已添加"}
}

// UpdateConnection 更新连接
func (cm *ConnectionManager) UpdateConnection(id, name, address, user string) OperationResult {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for i, c := range cm.connections {
		if c.ID == id {
			cm.connections[i].Name = name
			cm.connections[i].Address = address
			cm.connections[i].User = user
			if err := cm.saveConnections(); err != nil {
				return OperationResult{Success: false, Message: "保存失败: " + err.Error()}
			}
			return OperationResult{Success: true, Message: "连接已更新"}
		}
	}
	return OperationResult{Success: false, Message: "连接不存在"}
}

// DeleteConnection 删除连接
func (cm *ConnectionManager) DeleteConnection(id string) OperationResult {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for i, c := range cm.connections {
		if c.ID == id {
			cm.connections = append(cm.connections[:i], cm.connections[i+1:]...)
			if err := cm.saveConnections(); err != nil {
				return OperationResult{Success: false, Message: "保存失败: " + err.Error()}
			}
			return OperationResult{Success: true, Message: "连接已删除"}
		}
	}
	return OperationResult{Success: false, Message: "连接不存在"}
}

// GetConnection 获取单个连接
func (cm *ConnectionManager) GetConnection(id string) *ConnectionInfo {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for i := range cm.connections {
		if cm.connections[i].ID == id {
			return &cm.connections[i]
		}
	}
	return nil
}

func (cm *ConnectionManager) loadConnections() {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		cm.connections = []ConnectionInfo{}
		return
	}
	var conns []ConnectionInfo
	if err := json.Unmarshal(data, &conns); err != nil {
		cm.connections = []ConnectionInfo{}
		return
	}
	cm.connections = conns
}

func (cm *ConnectionManager) saveConnections() error {
	data, err := json.MarshalIndent(cm.connections, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cm.configPath, data, 0644)
}
