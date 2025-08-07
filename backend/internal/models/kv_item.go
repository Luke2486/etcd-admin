package models

import "time"

// KVItem 表示键值条目
type KVItem struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ConnectionID uint      `json:"connection_id" gorm:"not null;index"`
	Key          string    `json:"key" gorm:"not null;size:255"`
	Value        string    `json:"value" gorm:"type:text"` // 存储JSON字符串
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (KVItem) TableName() string {
	return "kv_items"
}
