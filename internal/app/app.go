package app

import (
	"fmt"

	"github.com/potom_pridumaem/config"
	restapi "github.com/potom_pridumaem/internal/controller"
	"github.com/potom_pridumaem/internal/repo/persistent"
	"github.com/potom_pridumaem/internal/usecase/property"
	"github.com/potom_pridumaem/internal/usecase/user"
	"github.com/potom_pridumaem/pkg/httpserver"
	"github.com/potom_pridumaem/pkg/logger"
	"github.com/potom_pridumaem/pkg/postgres"
)

type useCases struct {
	user     *user.UseCase
	property *property.UseCase
}

type repos struct {
	property *persistent.PropertyRepo
}

func initDeps(pg *postgres.Postgres) (useCases, repos) {
	userRepo := persistent.NewUserRepo(pg)
	propertyRepo := persistent.NewPropertyRepo(pg)

	return useCases{
			user:     user.New(userRepo),
			property: property.New(propertyRepo),
		}, repos{
			property: propertyRepo,
		}
}

func Run(cfg *config.Config) {
	lgr, err := logger.NewLogger(*cfg)
	if err != nil {
		panic(fmt.Sprintf("Init logger error: %s", err))
	}

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		panic(fmt.Sprintf("Postgres init error: %s", err))
	}
	defer pg.Close()

	uc, rp := initDeps(pg)

	httpServer := httpserver.NewServer(lgr.Logger)

	restapi.NewRouter(httpServer.App, cfg, uc.user, uc.property, rp.property, lgr.Logger)

	httpServer.Start()
	httpServer.WaitForShutdown(*lgr.Logger)
}
