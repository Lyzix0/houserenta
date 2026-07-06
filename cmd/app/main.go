package main

import (
	"github.com/potom_pridumaem/config"
	"github.com/potom_pridumaem/internal/app"
)

// @title           Potom Pridumaem API
// @version         1.0
// @description     API для управления объектами недвижимости и арендаторами
// @host            localhost:5050
// @BasePath        /v1
func main() {
	cfg := config.NewConfigMust()
	app.Run(&cfg)
}
