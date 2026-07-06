package main

import (
	"github.com/potom_pridumaem/config"
	_ "github.com/potom_pridumaem/docs"
	"github.com/potom_pridumaem/internal/app"
)

// @title           Potom Pridumaem API
// @version         1.0
// @description     API for managing rental properties, landlords, and tenants
// @host            localhost:5050
// @BasePath        /v1

// @securityDefinitions.apikey  CookieAuth
// @in                          cookie
// @name                        session_id
// @description                 Session cookie set on successful login
func main() {
	cfg := config.NewConfigMust()
	app.Run(&cfg)
}
