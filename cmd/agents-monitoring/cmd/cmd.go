package cmd

import (
	"context"
	"flag"
	"fmt"
	v1 "github.com/agents-monitoring/handlers/rest/v1"
	"github.com/agents-monitoring/model"
	"github.com/agents-monitoring/repository"
	"github.com/agents-monitoring/route"
	"github.com/agents-monitoring/service"
	echo "github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Execute() {
	e := echo.New()

	// flags
	// Note: These variables are here for simplicity.
	httpServerPort := flag.String("http-server-port", ":8080", "http server port")
	IpReceiverChannelBuffer := flag.Int("ip-receiver-buffer", 100, "buffer for the IpReceptor channel")

	// @TODO regenerate and move to secret manager
	IpInfoToken := flag.String("ipinfo-token", "e854f2dabdaf57", "ipinfo.io api token")
	dbUser := flag.String("db-user", "root", "ipinfo.io api token")
	dbSecret := flag.String("db-secret", "verysecuresecret", "ipinfo.io api token")
	dbHost := flag.String("db-host", "127.0.0.1", "ipinfo.io api token")
	dbName := flag.String("db-name", "agents", "ipinfo.io api token")
	flag.Parse()

	dbConnection, err := NewGormMySQLInstance(*dbHost, *dbName, *dbUser, *dbSecret)
	if err != nil {
		e.Logger.Fatal("failed to create database connection ", "error", err)
	}

	// Migrate the schema
	err = dbConnection.AutoMigrate(&model.AgentLog{}, model.AgentGeolocation{})
	if err != nil {
		e.Logger.Fatal("unable to migrate database ", "error", err)
	}

	agentRepository := repository.NewAgentRepository(dbConnection)

	ipInfoService := service.NewIpInfoService(*IpInfoToken, *IpReceiverChannelBuffer, agentRepository)
	go ipInfoService.Process(context.Background())
	router := &route.Router{
		Engine: e,
		V1: &v1.Handlers{
			IpInfoService:   ipInfoService,
			IpInfoKey:       IpInfoToken,
			AgentRepository: agentRepository,
		},
	}

	err = router.Execute()
	if err != nil {
		e.Logger.Fatal("Shutting down the server ", "error", err)
	}
	e.Logger.Fatal(e.Start(*httpServerPort))
}

func NewGormMySQLInstance(host string, databaseName string, username string, password string) (*gorm.DB, error) {
	var dbConnectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", username, password, host, databaseName)
	var db *gorm.DB
	var errConnecting error
	if db, errConnecting = gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{}); errConnecting != nil {
		return nil, errConnecting
	}

	return db, nil
}
