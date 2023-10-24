package v1

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) GetAgentLogs(ctx echo.Context) error {
	agentLogs, err := h.AgentRepository.GetAgentLogs(context.Background())
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, apiResponse{
		"agent_logs": agentLogs,
	})
}
