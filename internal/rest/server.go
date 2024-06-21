package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"social-network-otus/internal/session"

	"go.uber.org/fx"

	"github.com/gin-gonic/gin"

	"social-network-otus/internal/auth"
	"social-network-otus/internal/config"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/rest/handler"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/rest/router"
	"social-network-otus/internal/user"
)

type RestServer struct {
	Engine *gin.Engine
	Router *router.Router
	Logger logger.LoggerInterface
	Config *config.Config
}

type RestServerParams struct {
	fx.In

	Config      *config.Config
	Handler     *handler.RestHandler
	Logger      logger.LoggerInterface
	AuthService *auth.AuthService
	UserService *user.Service
	Session     *session.SessionStorage
	Response    *response.ResponseFactory
}

func NewRestServer(params RestServerParams) *RestServer {
	ginEngine := gin.Default()
	restRouter := router.NewRouter(ginEngine, params.Handler, params.Session, params.AuthService, params.UserService, params.Response)
	return &RestServer{
		Engine: ginEngine,
		Router: restRouter,
		Logger: params.Logger,
		Config: params.Config,
	}
}

func InitHooks(lc fx.Lifecycle, server *RestServer) {
	var srv *http.Server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			server.Logger.Info(fmt.Sprintf("Starting HTTP server at %s", server.Config.HttpServer), map[string]interface{}{})
			go func() {
				srv = &http.Server{
					Addr:    server.Config.HttpServer,
					Handler: server.Engine,
				}
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					server.Logger.Fatal("start listen and serve error", err, map[string]interface{}{})
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
