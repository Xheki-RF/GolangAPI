package assembly

import (
	"net/http"
	"social-network/conf"
	"social-network/controller"
	"social-network/repository"
	"social-network/routes"
	"social-network/service"

	"github.com/txix-open/grmq/consumer"
	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/grmqx"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"github.com/txix-open/isp-kit/grpc/endpoint/grpclog"
	httpEndpoint "github.com/txix-open/isp-kit/http/endpoint"
	"github.com/txix-open/isp-kit/http/endpoint/httplog"
	"github.com/txix-open/isp-kit/log"
)

type DB interface {
	db.DB
	db.Transactional
}

type Locator struct {
	db     DB
	logger log.Logger
}

type LocatorConfig struct {
	HttpHandler http.Handler
	GrpcHandler *grpc.Mux
	RmqHandler  consumer.Consumer
}

func NewLocator(db DB, logger log.Logger) Locator {
	return Locator{
		db:     db,
		logger: logger,
	}
}

func (l Locator) Config(conf conf.Remote) LocatorConfig {
	userRepo := repository.NewUser(l.logger, l.db, l.db)
	userService := service.NewUser(userRepo)
	userController := controller.NewUser(userService)

	c := routes.Controllers{
		User: *userController,
	}

	mapper := endpoint.DefaultWrapper(l.logger, grpclog.Log(l.logger, true))
	grpcHandler := routes.Handler(mapper, c)

	wrapper := httpEndpoint.DefaultWrapper(l.logger, httplog.Log(l.logger, true))
	httpHandler := routes.HttpHandler(wrapper, c)

	handler := grmqx.NewResultHandler(l.logger, nil)
	rmqHandler := conf.Consumer.Config.DefaultConsumer(handler, grmqx.ConsumerLog(l.logger, true))

	return LocatorConfig{
		HttpHandler: httpHandler,
		GrpcHandler: grpcHandler,
		RmqHandler:  rmqHandler,
	}
}
