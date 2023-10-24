package model

import "time"

type AgentLog struct {
	IpAddress string     `gorm:"column:ip_address" json:"ip_address"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
}
