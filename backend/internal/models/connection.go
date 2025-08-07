package models

import (
	"time"

	"gorm.io/gorm"
)

// Connection 表示etcd连接配置
type Connection struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null;size:100;uniqueIndex"`
	Endpoints   string         `json:"endpoints" gorm:"not null;type:text"` // JSON数组字符串，如["localhost:2379"]
	Username    string         `json:"username" gorm:"size:100"`
	Password    string         `json:"password" gorm:"size:255"`
	TLSEnabled  bool           `json:"tls_enabled" gorm:"default:false"`
	CertFile    string         `json:"cert_file" gorm:"size:255"`
	KeyFile     string         `json:"key_file" gorm:"size:255"`
	CAFile      string         `json:"ca_file" gorm:"size:255"`
	Description string         `json:"description" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	IsReadOnly  bool           `json:"is_readonly" gorm:"column:is_readonly;default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	KVItems []KVItem `json:"kv_items,omitempty" gorm:"foreignKey:ConnectionID"`
}

// TableName 指定表名
func (Connection) TableName() string {
	return "connections"
}
