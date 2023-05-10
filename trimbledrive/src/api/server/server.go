package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	adapter "github.com/gwatts/gin-adapter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/configuration"
	"trimble.com/tdrive/api/config"
	"trimble.com/tdrive/api/middleware"
	"trimble.com/tdrive/api/middleware/compress"
	"trimble.com/tdrive/api/router"
	"trimble.com/tdrive/api/utils"
)

func Initialize(
	tDriveConfig *config.TDriveConfig,
	logger *zap.Logger,
	requestIdMiddleware middleware.RequestIdMiddleware,
	accessLoggerMiddleware middleware.AccessLogMiddleware,
	customRecoveryMiddleware middleware.CustomRecoveryMiddleware,
	authMiddleware middleware.AuthMiddleware,
	corsMiddleware middleware.CorsMiddleware,
	cacheControlMiddleware middleware.CacheControlMiddleware,
	decompressionMiddleware middleware.DecompressionMiddleware,
	routers router.Routers,
	mongoClient *mongo.Client,
) {
	// Create new blank gin engine instance
	ginEngine := gin.New()
	// Registering custom validators
	registerValidators()
	// Add required middleware
	ginEngine.Use(requestIdMiddleware.HandleRequestId())
	ginEngine.Use(accessLoggerMiddleware.HandleAccessLog())
	ginEngine.Use(corsMiddleware.HandleCors())
	ginEngine.Use(cacheControlMiddleware.HandleCacheControl())

	ginEngine.Use(decompressionMiddleware.HandleDecompression())
	if tDriveConfig.Compress.MinResponseBodyBytes >= 0 {
		extHandler, wrapper := adapter.New()
		compressMiddleware, err := compress.New(logger, tDriveConfig.Compress, extHandler)
		if err == nil {
			logger.Debug("Installing compress middleware")
			ginEngine.Use(wrapper(compressMiddleware))
		}
	}

	ginEngine.Use(customRecoveryMiddleware.HandlerCustomRecovery())

	ginEngine.NoRoute(func(c *gin.Context) {
		utils.HandleApiError(c, api_error.InvalidUrl)
	})

	ginEngine.HandleMethodNotAllowed = true
	ginEngine.NoMethod(func(c *gin.Context) {
		utils.HandleApiError(c, api_error.MethodNotAllowed)
	})

	// Add ginEngine config
	appV1 := ginEngine.Group("v1")
	routers.Register(appV1)

	// Listen Server on specified port
	port := ":" + strconv.Itoa(tDriveConfig.Server.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: ginEngine,
	}

	// Below is to handle Graceful Shutdown

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	log.Println("Initializing server and listening on port: " + port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error initializing server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gin server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// End JWKS refresh go routine
	authMiddleware.JWKS.EndBackground()

	// Disconnect mongo client
	configuration.CloseMongoConnection(mongoClient)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("gin Server forced to shutdown: ", zap.Error(err))
	}
	log.Println("gin Server exiting")
	// Below is to trigger Fx App shutdown
	os.Exit(1)
}

// registerValidators Registers custom validators
func registerValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		urlSafeErr := v.RegisterValidation("UrlSafe", utils.UrlSafe)
		if urlSafeErr != nil {
			panic("Error on registering validation")
		}
	}
}
