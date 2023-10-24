package v1

import (
	"net/http"

	"github.com/agents-monitoring/repository"
	"github.com/agents-monitoring/service"
	echo "github.com/labstack/echo/v4"
)

const badRequestErrorMessage string = "invalid parameters in the request - Error: %s"

type Handlers struct {
	IpInfoService   service.IpInfoService
	IpInfoKey       *string
	AgentRepository repository.AgentRepository
}

type apiResponse map[string]interface{}

func (h *Handlers) HealthCheck(ctx echo.Context) error {
	res := apiResponse{"status": "OK"}
	return ctx.JSON(http.StatusOK, res)
}
