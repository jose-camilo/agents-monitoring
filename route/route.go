package route

import (
	"errors"
	v1 "github.com/agents-monitoring/handlers/rest/v1"
	echo "github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

// Router object will contain handlers by version
type Router struct {
	Engine *echo.Echo
	V1     *v1.Handlers
}

func (r *Router) Execute() error {

	if r.Engine == nil {
		return errors.New("router.Engine top level echo framework instance is nil, cannot be nil")
	}

	if r.V1 == nil {
		return errors.New("router.V1 needs to be set, cannot be nil")
	}

	r.Engine.Pre(echomw.RemoveTrailingSlash())
	v1Route := r.Engine.Group("/v1")
	v1Route.GET("/health-check", r.V1.HealthCheck)
	v1Route.POST("/audit-log", r.V1.CreateAgentLog)
	v1Route.GET("/audit-logs", r.V1.GetAgentLogs)
	v1Route.GET("/agent/:ipaddress", r.V1.GetAgentByIpAddress)

	return nil
}
