package api

import (
	"context"
	"net/http"
	"time"

	"github.com/eztwokey/l3-3/internal/config"
	"github.com/eztwokey/l3-3/internal/logic"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

type Api struct {
	server *http.Server
	engine *gin.Engine
	logic  *logic.Logic
	logger logger.Logger
}

func New(cfg *config.Config, logic *logic.Logic, logger logger.Logger) *Api {
	gin.SetMode(cfg.Api.GinMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	if cfg.Api.GinMode == gin.DebugMode {
		engine.Use(gin.Logger())
	}

	server := &http.Server{
		Addr:         cfg.Api.Addr,
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Api.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Api.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Api.IdleTimeout) * time.Second,
	}

	api := &Api{
		server: server,
		engine: engine,
		logic:  logic,
		logger: logger,
	}
	api.registerRoutes()

	return api
}

func (a *Api) Run() error {
	return a.server.ListenAndServe()
}

func (a *Api) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *Api) registerRoutes() {
	a.engine.StaticFile("/", "./web/index.html")
	a.engine.POST("/comments", a.createComment)
	a.engine.GET("/comments", a.listComments)
	a.engine.GET("/comments/:id", a.getTree)
	a.engine.DELETE("/comments/:id", a.deleteComment)
}
