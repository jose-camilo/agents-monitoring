package v1

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/agents-monitoring/model"
	echo "github.com/labstack/echo/v4"
)

func (h *Handlers) CreateAgentLog(ctx echo.Context) error {
	agentLog := &model.AgentLog{}
	if err := ctx.Bind(agentLog); err != nil {
		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf(badRequestErrorMessage, err.Error()),
		)
	}

	timeNow := time.Now()
	agentLog.CreatedAt = &timeNow

	if net.ParseIP(agentLog.IpAddress) == nil {
		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf("invalid ip address in the request"),
		)
	}

	err := h.AgentRepository.CreateAgentLog(context.Background(), agentLog)
	if err != nil {
		ctx.Logger().Error("handler.CreateAgentLog: unable to create agent log", "error", err)
	}

	h.IpInfoService.Add(&agentLog.IpAddress)

	return ctx.NoContent(http.StatusNoContent)
}
