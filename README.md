# Agents Monitoring
##Requirements 
* Docker
* Golang

## Run MySQL on Docker
```bash
docker run --name agents -e MYSQL_ROOT_PASSWORD=verysecuresecret -e MYSQL_DATABASE=agents -p 3306:3306 -d mysql:8.0
```

## Build the application
```bash
go build -o app ./cmd/agents-monitoring/
```

## Run the application
```bash
./app
```

## Endpoints
#### http://127.0.0.1:8080/v1/health-check
Request GET
#### http://127.0.0.1:8080/v1/audit-log
Request POST
```json
{
  "ip_address": "{ip_address}"
}

```
#### http://127.0.0.1:8080/v1/audit-logs
Request GET
#### http://127.0.0.1:8080/v1/agent/:ipaddress
Request GET