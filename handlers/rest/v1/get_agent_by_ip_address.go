package v1

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func (h *Handlers) GetAgentByIpAddress(ctx echo.Context) error {
	ipAddress := ctx.Param("ipaddress")
	if ipAddress == "" {
		return errors.New("missing field on request")
	}

	if net.ParseIP(ipAddress) == nil {
		return ctx.String(
			http.StatusBadRequest,
			fmt.Sprintf("invalid ip address in the request"),
		)
	}

	agentGeolocation, err := h.AgentRepository.GetAgentGeolocation(context.Background(), ipAddress)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, apiResponse{
		"agent_geolocation": agentGeolocation,
	})
}
