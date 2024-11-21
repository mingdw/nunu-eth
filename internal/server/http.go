package server

import (
	"nunu-eth/docs"
	"nunu-eth/internal/handler"
	"nunu-eth/internal/middleware"
	"nunu-eth/pkg/jwt"
	"nunu-eth/pkg/log"
	nunuhttp "nunu-eth/pkg/server/http"
	"nunu-eth/static"
	"nunu-eth/web"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	commonHandler *handler.CommonHandler,
) *nunuhttp.Server {
	gin.SetMode(gin.DebugMode)
	s := nunuhttp.NewServer(
		gin.Default(),
		logger,
		nunuhttp.WithServerHost(conf.GetString("http.host")),
		nunuhttp.WithServerPort(conf.GetInt("http.port")),
	)
	s.Use(static.Serve("/", static.EmbedFolder(web.HtmlsFs, ".")))

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	api := s.Group("/api")
	{
		api.POST("/connectTest", commonHandler.TestConnectClient)

		api.GET("/accountFormat", commonHandler.AccountFormt)

		api.POST("/accountBalance", commonHandler.AccountBalance)

		api.POST("/blockQuery", commonHandler.BlockQuery)

		api.POST("/transactionQuery", commonHandler.TransactionQuery)
	}
	v1 := s.Group("/v1")
	{
		// No route group has permission
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/register", userHandler.Register)
			noAuthRouter.POST("/login", userHandler.Login)
			noAuthRouter.GET("/common", commonHandler.Test)
		}
		// Non-strict permission routing group
		noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			noStrictAuthRouter.GET("/user", userHandler.GetProfile)
		}

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
		}

	}
	return s
}
