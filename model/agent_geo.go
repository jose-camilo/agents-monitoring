package model

import "time"

type AgentGeolocation struct {
	Ip        string     `gorm:"column:ip" json:"ip"`
	Hostname  string     `gorm:"column:hostname" json:"hostname"`
	City      string     `gorm:"column:city" json:"city"`
	Region    string     `gorm:"column:region" json:"region"`
	Country   string     `gorm:"column:country" json:"country"`
	Loc       string     `gorm:"column:loc" json:"loc"`
	Org       string     `gorm:"column:-" json:"org"`
	Isp       string     `gorm:"column:isp" json:"Isp"`
	Asn       string     `gorm:"column:asn" json:"Asn"`
	Postal    string     `gorm:"column:postal" json:"postal"`
	Timezone  string     `gorm:"column:timezone" json:"timezone"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
}
