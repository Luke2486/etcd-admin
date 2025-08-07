package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"etcd-admin-backend/internal/models"
)

// EtcdService etcd客户端服务
type EtcdService struct {
	clients map[uint]*clientv3.Client // connection_id -> client
}

// NewEtcdService 创建etcd服务实例
func NewEtcdService() *EtcdService {
	return &EtcdService{
		clients: make(map[uint]*clientv3.Client),
	}
}

// GetClient 获取或创建etcd客户端
func (s *EtcdService) GetClient(conn *models.Connection) (*clientv3.Client, error) {
	// 检查是否已存在客户端
	if client, exists := s.clients[conn.ID]; exists {
		return client, nil
	}

	// 解析endpoints
	var endpoints []string
	if err := json.Unmarshal([]byte(conn.Endpoints), &endpoints); err != nil {
		// 如果不是JSON格式，尝试按逗号分割
		endpoints = strings.Split(conn.Endpoints, ",")
		for i, endpoint := range endpoints {
			endpoints[i] = strings.TrimSpace(endpoint)
		}
	}

	// 處理localhost地址，在Docker容器中使用extra_hosts配置
	// Docker的extra_hosts已經將localhost映射到正確的主機地址
	// 無需手動替換，保持原始的localhost地址

	// 创建etcd客户端配置
	config := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	// 如果有用户名和密码，设置认证
	if conn.Username != "" && conn.Password != "" {
		config.Username = conn.Username
		config.Password = conn.Password
	}

	// 创建客户端
	client, err := clientv3.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = client.Status(ctx, endpoints[0])
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	// 缓存客户端
	s.clients[conn.ID] = client
	return client, nil
}

// CloseClient 关闭特定连接的客户端
func (s *EtcdService) CloseClient(connectionID uint) {
	if client, exists := s.clients[connectionID]; exists {
		client.Close()
		delete(s.clients, connectionID)
	}
}

// CloseAll 关闭所有客户端
func (s *EtcdService) CloseAll() {
	for id, client := range s.clients {
		client.Close()
		delete(s.clients, id)
	}
}

// ListKeys 列出所有keys
func (s *EtcdService) ListKeys(conn *models.Connection, prefix string) ([]string, error) {
	client, err := s.GetClient(conn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	keys := make([]string, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		keys[i] = string(kv.Key)
	}

	return keys, nil
}

// GetValue 获取键值
func (s *EtcdService) GetValue(conn *models.Connection, key string) (string, error) {
	client, err := s.GetClient(conn)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get value: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return string(resp.Kvs[0].Value), nil
}

// SetValue 设置键值
func (s *EtcdService) SetValue(conn *models.Connection, key, value string) error {
	client, err := s.GetClient(conn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

// DeleteKey 删除键
func (s *EtcdService) DeleteKey(conn *models.Connection, key string) error {
	client, err := s.GetClient(conn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// GetAllKV 获取所有键值对
func (s *EtcdService) GetAllKV(conn *models.Connection, prefix string) (map[string]string, error) {
	client, err := s.GetClient(conn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get all KV: %w", err)
	}

	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		result[string(kv.Key)] = string(kv.Value)
	}

	return result, nil
}

// TestConnection 测试连接
func (s *EtcdService) TestConnection(conn *models.Connection) error {
	client, err := s.GetClient(conn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 解析endpoints获取第一个进行测试
	var endpoints []string
	if err := json.Unmarshal([]byte(conn.Endpoints), &endpoints); err != nil {
		endpoints = strings.Split(conn.Endpoints, ",")
		for i, endpoint := range endpoints {
			endpoints[i] = strings.TrimSpace(endpoint)
		}
	}

	if len(endpoints) == 0 {
		return fmt.Errorf("no endpoints configured")
	}

	_, err = client.Status(ctx, endpoints[0])
	return err
}
