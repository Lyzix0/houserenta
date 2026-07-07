package app

import (
	"fmt"

	"github.com/potom_pridumaem/config"
	restapi "github.com/potom_pridumaem/internal/controller"
	"github.com/potom_pridumaem/internal/repo/persistent"
	"github.com/potom_pridumaem/internal/usecase/billing"
	"github.com/potom_pridumaem/internal/usecase/property"
	"github.com/potom_pridumaem/internal/usecase/user"
	"github.com/potom_pridumaem/pkg/httpserver"
	"github.com/potom_pridumaem/pkg/logger"
	"github.com/potom_pridumaem/pkg/postgres"
)

type useCases struct {
	user     *user.UseCase
	property *property.UseCase
	billing  *billing.UseCase
}

func initUseCases(pg *postgres.Postgres) useCases {
	userRepo := persistent.NewUserRepo(pg)
	propertyRepo := persistent.NewPropertyRepo(pg)
	leaseRepo := persistent.NewLeaseRepo(pg)
	readingRepo := persistent.NewReadingRepo(pg)
	billRepo := persistent.NewBillRepo(pg)
	customNextItemRepo := persistent.NewCustomNextItemRepo(pg)

	return useCases{
		user:     user.New(userRepo, leaseRepo),
		property: property.New(propertyRepo, leaseRepo, readingRepo, billRepo, customNextItemRepo, userRepo),
		billing:  billing.New(leaseRepo, propertyRepo, billRepo, readingRepo, customNextItemRepo),
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

	uc := initUseCases(pg)

	httpServer := httpserver.NewServer(lgr.Logger, httpserver.AllowOrigins(cfg.HTTP.AllowedOrigins))

	restapi.NewRouter(httpServer.App, cfg, uc.user, uc.property, uc.billing, lgr.Logger)

	httpServer.Start()
	httpServer.WaitForShutdown(*lgr.Logger)
}
