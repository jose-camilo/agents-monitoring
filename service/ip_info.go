package service

//go:generate mockgen -destination=mock/mock_ip_info.go -package=mock github.com/jose-camilo/agents-monitoring/service IpInfoService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/agents-monitoring/model"
	"github.com/agents-monitoring/repository"
)

const GetIpInfoUrl string = "https://ipinfo.io/%s?token=%s"

type IpInfoService interface {
	Process(ctx context.Context)
	GetData(_ context.Context, ipAddressMessage *IpAddressQueueMessage) (*model.AgentGeolocation, error)
	Add(ipAddress *string)
}

type IpAddressQueueMessage struct {
	IpAddress    *string
	TimesRetried int
}

type IpInfoServiceImpl struct {
	token           string
	IpAddressQueue  chan *IpAddressQueueMessage
	AgentRepository repository.AgentRepository
}

func NewIpInfoService(token string, buffer int, agentRepository repository.AgentRepository) IpInfoService {
	ipInfoChan := make(chan *IpAddressQueueMessage, buffer)
	return &IpInfoServiceImpl{
		token:           token,
		IpAddressQueue:  ipInfoChan,
		AgentRepository: agentRepository,
	}
}

func (i *IpInfoServiceImpl) Add(ipAddress *string) {
	i.IpAddressQueue <- &IpAddressQueueMessage{
		IpAddress: ipAddress,
	}
}

func (i *IpInfoServiceImpl) Process(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR: ip-info queue: ip info channel failure")
		}
	}()

	for {
		ctxProcess := context.Background()
		ipAddressQueueMessage := <-i.IpAddressQueue

		if ipAddressQueueMessage.TimesRetried > 3 {
			// @TODO store this ip address somewhere other than the logs as it's a GDPR/CCPA topic.
			fmt.Println("ERROR: ip-info queue: ip address data gathering failure after 3 attempts")
			continue
		}

		if ipAddressQueueMessage.IpAddress == nil {
			fmt.Println("ERROR: ip-info queue: empty queue message")
			continue
		}

		// TODO add ratelimit to this logic
		agentGeolocation, err := i.GetData(ctxProcess, ipAddressQueueMessage)
		if err != nil {
			// @TODO check if failure message contains PII
			fmt.Println("ERROR: ip-info queue: failed to retrieve data for ip address" + err.Error())
			ipAddressQueueMessage.TimesRetried++
			i.IpAddressQueue <- ipAddressQueueMessage
			continue
		}

		_, err = i.AgentRepository.GetAgentGeolocation(ctxProcess, agentGeolocation.Ip)
		if err != nil && err.Error() == "not found" {
			err = i.AgentRepository.CreateAgentGeolocation(ctxProcess, agentGeolocation)
			if err != nil {
				fmt.Println("ERROR: ip-info queue: failed to store agent geolocation data on db" + err.Error())
				ipAddressQueueMessage.TimesRetried++
				i.IpAddressQueue <- ipAddressQueueMessage
			}
		}
	}
}

func (i *IpInfoServiceImpl) GetData(_ context.Context, ipAddressMessage *IpAddressQueueMessage) (*model.AgentGeolocation, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(GetIpInfoUrl, *ipAddressMessage.IpAddress, i.token), nil)
	if err != nil {
		return &model.AgentGeolocation{}, err
	}
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("ERROR: Failed to close body for ip info data request")
		}
	}(resp.Body)
	if err != nil {
		return &model.AgentGeolocation{}, err
	}

	if resp.StatusCode > 299 {
		return &model.AgentGeolocation{}, errors.New(fmt.Sprintf("request failed with status code %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &model.AgentGeolocation{}, err
	}

	ipInfoResponse := &model.AgentGeolocation{}
	err = json.Unmarshal(body, ipInfoResponse)
	if err != nil {
		return &model.AgentGeolocation{}, err
	}
	return ipInfoResponse, nil
}
