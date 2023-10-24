package repository

import (
	"context"
	"errors"
	"github.com/agents-monitoring/model"
	"gorm.io/gorm"
)

type AgentRepository interface {
	// Agent Log
	CreateAgentLog(ctx context.Context, c *model.AgentLog) error
	GetAgentLogs(ctx context.Context) ([]model.AgentLog, error)

	// Agent Geolocation
	CreateAgentGeolocation(ctx context.Context, c *model.AgentGeolocation) error
	GetAgentGeolocation(ctx context.Context, ipAddress string) (*model.AgentGeolocation, error)
	UpdateAgentGeolocation(ctx context.Context, agentGeolocation *model.AgentGeolocation) error
}

type AgentRepositoryImpl struct {
	DB *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &AgentRepositoryImpl{
		DB: db,
	}
}

func (i *AgentRepositoryImpl) CreateAgentLog(ctx context.Context, c *model.AgentLog) error {
	return i.DB.Create(&c).Error
}

func (i *AgentRepositoryImpl) CreateAgentGeolocation(ctx context.Context, c *model.AgentGeolocation) error {
	return i.DB.Create(&c).Error
}

func (i *AgentRepositoryImpl) GetAgentGeolocation(ctx context.Context, ipAddress string) (*model.AgentGeolocation, error) {
	agentGeolocation := model.AgentGeolocation{}

	if err := i.DB.First(&agentGeolocation, "ip = ?", ipAddress).Error; gorm.ErrRecordNotFound == err {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}

	return &agentGeolocation, nil
}

func (i *AgentRepositoryImpl) UpdateAgentGeolocation(ctx context.Context, agentGeolocation *model.AgentGeolocation) error {
	return i.DB.Model(&model.AgentGeolocation{Ip: agentGeolocation.Ip}).Updates(model.AgentGeolocation{
		Ip:        agentGeolocation.Ip,
		Hostname:  agentGeolocation.Hostname,
		City:      agentGeolocation.City,
		Region:    agentGeolocation.Region,
		Country:   agentGeolocation.Country,
		Loc:       agentGeolocation.Loc,
		Isp:       agentGeolocation.Isp,
		Asn:       agentGeolocation.Asn,
		Postal:    agentGeolocation.Postal,
		Timezone:  agentGeolocation.Timezone,
		CreatedAt: agentGeolocation.CreatedAt,
	}).Error
}

func (i *AgentRepositoryImpl) GetAgentLogs(ctx context.Context) ([]model.AgentLog, error) {
	agentLogs := []model.AgentLog{}
	i.DB.Raw("SELECT * FROM agent_logs LIMIT 1000").Scan(&agentLogs)
	return agentLogs, nil
}
