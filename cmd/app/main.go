package main

import (
	"github.com/potom_pridumaem/config"
	"github.com/potom_pridumaem/internal/app"
)

func main() {
	cfg := config.NewConfigMust()
	app.Run(&cfg)
}
